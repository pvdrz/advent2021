use crate::counting_iter::CountingIter;

const fn extract_usize<'a>(
    mut bits: CountingIter<'a>,
    len_bits: usize,
) -> (usize, CountingIter<'a>) {
    let mut i = 0;
    let mut int = 0;
    while i < len_bits {
        int *= 2;

        let (new_bits, bit) = bits.next();
        bits = new_bits;

        if bit {
            int += 1;
        }

        i += 1;
    }

    return (int, bits);
}

macro_rules! eval_bin {
    ($name:ident, $op:tt) => {
        const fn $name(
            mut bits: CountingIter<'_>,
            version_acc: usize,
        ) -> (usize, usize, CountingIter<'_>) {
            let (new_bits, length_type_id) = bits.next();
            bits = new_bits;

            if length_type_id {
                let (subpacket_count, new_bits) = extract_usize(bits, 11);
                bits = new_bits;

                if subpacket_count != 2 {
                    panic!(concat!("Expected a subpacket count of two for ", stringify!($name)));
                }

                let (version_a, res_a, bits) = eval_inner(bits);
                let (version_b, res_b, bits) = eval_inner(bits);

                (
                    version_acc + version_a + version_b,
                    (res_a $op res_b) as usize,
                    bits,
                )
            } else {
                let (total_length, new_bits) = extract_usize(bits, 15);
                bits = new_bits;

                let currently_used = bits.calls();

                let (version_a, res_a, bits) = eval_inner(bits);
                let (version_b, res_b, bits) = eval_inner(bits);

                if bits.calls() != currently_used + total_length {
                    panic!("Invalid bit count");
                }
                (
                    version_acc + version_a + version_b,
                    (res_a $op res_b) as usize,
                    bits,
                )
            }
        }
    };
}

macro_rules! eval_fold {
    ($name:ident, $op:ident) => {
        const fn $name(
            mut bits: CountingIter<'_>,
            mut version_acc: usize,
            mut acc: usize,
        ) -> (usize, usize, CountingIter<'_>) {
            let (new_bits, length_type_id) = bits.next();
            bits = new_bits;

            if length_type_id {
                let (subpacket_count, new_bits) = extract_usize(bits, 11);
                bits = new_bits;

                let mut i = 0;
                while i < subpacket_count {
                    let (version, res, new_bits) = eval_inner(bits);
                    bits = new_bits;

                    version_acc += version;
                    acc = $op(acc, res);

                    i += 1;
                }
            } else {
                let (total_length, new_bits) = extract_usize(bits, 15);
                bits = new_bits;

                let currently_used = bits.calls();

                while bits.calls() < currently_used + total_length {
                    let (version, res, new_bits) = eval_inner(bits);
                    bits = new_bits;

                    version_acc += version;
                    acc = $op(acc, res);
                }
            }

            (version_acc, acc, bits)
        }
    };
}

pub(crate) const fn eval(s: &[u8]) -> (usize, usize) {
    let iter = CountingIter::new(s);
    let (version_sum, res, _) = eval_inner(iter);
    (version_sum, res)
}

const fn eval_inner(mut bits: CountingIter<'_>) -> (usize, usize, CountingIter<'_>) {
    let (version, new_bits) = extract_usize(bits, 3);
    bits = new_bits;

    let (kind, new_bits) = extract_usize(bits, 3);
    bits = new_bits;

    const fn add(x: usize, y: usize) -> usize {
        x + y
    }
    const fn mul(x: usize, y: usize) -> usize {
        x * y
    }
    const fn min(x: usize, y: usize) -> usize {
        if x < y {
            x
        } else {
            y
        }
    }
    const fn max(x: usize, y: usize) -> usize {
        if x > y {
            x
        } else {
            y
        }
    }

    eval_fold!(eval_add, add);
    eval_fold!(eval_prod, mul);
    eval_fold!(eval_min, min);
    eval_fold!(eval_max, max);
    eval_bin!(eval_gt, >);
    eval_bin!(eval_lt, <);
    eval_bin!(eval_eq, ==);

    match kind {
        0 => eval_add(bits, version, 0),
        1 => eval_prod(bits, version, 1),
        2 => eval_min(bits, version, usize::MAX),
        3 => eval_max(bits, version, usize::MIN),
        4 => eval_lit(bits, version),
        5 => eval_gt(bits, version),
        6 => eval_lt(bits, version),
        7 => eval_eq(bits, version),
        _ => panic!("invalid payload kind"),
    }
}

const fn eval_lit(mut bits: CountingIter<'_>, version: usize) -> (usize, usize, CountingIter<'_>) {
    let mut literal = 0;

    loop {
        let (new_bits, more_packets) = bits.next();
        bits = new_bits;

        let (val, new_bits) = extract_usize(bits, 4);
        bits = new_bits;

        literal *= 2usize.pow(4);
        literal += val;

        if !more_packets {
            break;
        }
    }

    (version, literal, bits)
}
