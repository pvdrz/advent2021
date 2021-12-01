use anyhow::{Context, Error};

use std::fs::read_to_string;

#[inline(never)]
fn count_increases<const N: usize>(slice: &[u32]) -> usize {
    (0..slice.len() - N)
        .filter(|&i| slice[i + N] > slice[i])
        .count()
}

fn main() -> Result<(), Error> {
    let input = read_to_string("./input")
        .context("could not read input file")?
        .lines()
        .map(|line| line.parse().context("could not parse input line"))
        .collect::<Result<Vec<_>, _>>()?;

    println!("Part 1: {}", count_increases::<1>(&input));
    println!("Part 2: {}", count_increases::<3>(&input));

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test1() {
        let input = [199, 200, 208, 210, 200, 207, 240, 269, 260, 263];

        assert_eq!(7, count_increases::<1>(&input));
    }

    #[test]
    fn test2() {
        let input = [199, 200, 208, 210, 200, 207, 240, 269, 260, 263];

        assert_eq!(5, count_increases::<3>(&input));
    }
}
