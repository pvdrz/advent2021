use core::ops::{Add, Mul};

use crate::counting_iter::{CountingIter, CountingIterator};

type Error = &'static str;

pub fn evall(s: &str) -> Result<(usize, usize), Error> {
    eval_inner(&mut CountingIter::new(
        s.chars().flat_map(crate::hexa::to_binary),
    ))
}

pub fn eval_inner(bits: &mut impl CountingIterator<bool>) -> Result<(usize, usize), Error> {
    let mut version = 0;
    for _ in 0..3 {
        version *= 2;
        if bits.next().ok_or("Missing bit in version")? {
            version += 1;
        }
    }

    let mut kind = 0;
    for _ in 0..3 {
        kind *= 2;
        if bits.next().ok_or("Missing bit in type id")? {
            kind += 1;
        }
    }

    match kind {
        // Sum
        0 => eval_fold(bits, version, 0, Add::add),
        // Prod
        1 => eval_fold(bits, version, 1, Mul::mul),
        // Min
        2 => eval_fold(bits, version, usize::MAX, Ord::min),
        // Max
        3 => eval_fold(bits, version, usize::MIN, Ord::max),
        // Lit
        4 => return Ok((version, eval_lit(bits)?)),
        // Gt
        5 => eval_bin(bits, version, |x, y| x > y),
        // Lt
        6 => eval_bin(bits, version, |x, y| x < y),
        // Eq
        7 => eval_bin(bits, version, |x, y| x == y),
        _ => panic!("invalid payload kind"),
    }
}
fn eval_bin(
    bits: &mut impl CountingIterator<bool>,
    version_acc: usize,
    bin_op: fn(usize, usize) -> bool,
) -> Result<(usize, usize), Error> {
    let length_type_id = bits.next().ok_or("Missing length_type_id")?;

    if length_type_id {
        let mut subpacket_count = 0;
        for _ in 0..11 {
            subpacket_count *= 2;
            if bits.next().ok_or("Missing bit in subpacket count")? {
                subpacket_count += 1;
            }
        }

        if subpacket_count != 2 {
            return Err("Expected a subpacket count of two");
        }
        let (version_a, res_a) = eval_inner(bits)?;
        let (version_b, res_b) = eval_inner(bits)?;

        Ok((
            version_acc + version_a + version_b,
            bin_op(res_a, res_b) as usize,
        ))
    } else {
        let mut total_length = 0;
        for _ in 0..15 {
            total_length *= 2;
            if bits.next().ok_or("Missing bit in total length")? {
                total_length += 1;
            }
        }
        let currently_used = bits.calls();

        let (version_a, res_a) = eval_inner(bits)?;
        let (version_b, res_b) = eval_inner(bits)?;

        if bits.calls() != currently_used + total_length {
            return Err("Invalid bit count");
        }
        Ok((
            version_acc + version_a + version_b,
            bin_op(res_a, res_b) as usize,
        ))
    }
}

fn eval_fold(
    bits: &mut impl CountingIterator<bool>,
    mut version_acc: usize,
    mut acc: usize,
    fold: fn(usize, usize) -> usize,
) -> Result<(usize, usize), Error> {
    let length_type_id = bits.next().ok_or("Missing length_type_id")?;

    if length_type_id {
        let mut subpacket_count = 0;
        for _ in 0..11 {
            subpacket_count *= 2;
            if bits.next().ok_or("Missing bit in subpacket count")? {
                subpacket_count += 1;
            }
        }
        for _ in 0..subpacket_count {
            let (version, res) = eval_inner(bits)?;
            version_acc += version;
            acc = fold(acc, res);
        }
    } else {
        let mut total_length = 0;
        for _ in 0..15 {
            total_length *= 2;
            if bits.next().ok_or("Missing bit in total length")? {
                total_length += 1;
            }
        }
        let currently_used = bits.calls();

        while bits.calls() < currently_used + total_length {
            let (version, res) = eval_inner(bits)?;
            version_acc += version;
            acc = fold(acc, res);
        }
    }

    Ok((version_acc, acc))
}

pub fn eval_lit(bits: &mut impl CountingIterator<bool>) -> Result<usize, Error> {
    let mut literal = 0;
    while let Some(more_packets) = bits.next() {
        for _ in 0..4 {
            literal *= 2;
            if bits.next().ok_or("Missing bit in payload group")? {
                literal += 1;
            }
        }
        if !more_packets {
            break;
        }
    }
    Ok(literal)
}
