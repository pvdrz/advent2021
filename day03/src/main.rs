use anyhow::Error;

/// Counts the number of one-digits in `numbers` for each position.
fn count_ones(numbers: &[u64], digits: usize) -> Vec<usize> {
    // This array will contain the number of ocurrences of the digit one in each position.
    let mut one_counts = vec![0; digits];

    // Iterate over every number in the input.
    for number in numbers {
        // Iterate over each position and update the number of one-ocurrences if required.
        for (pos, one_count) in one_counts.iter_mut().enumerate() {
            // If the digit in the position `pos` is `1`, increase the count for that position.
            if (number >> pos) & 1 == 1 {
                *one_count += 1;
            }
        }
    }

    one_counts
}

/// Compute the gamma and epsilon rate from the count of one-digits and the total amount of
/// numbers.
fn compute_rates(one_counts: &[usize], numbers_len: usize) -> (u64, u64) {
    // The gamma rate is composed by the most common digits of each position. All digits are set to
    // zero.
    let mut gamma_rate: u64 = 0;

    // Iterate over each position and get the number of one-ocurrences.
    for (pos, one_count) in one_counts.iter().enumerate() {
        // If more than half of the numbers had a `1` in the current position, the digit of the
        // gamma rate in this position is `1`. Otherwise we left it be zero.
        if 2 * one_count >= numbers_len {
            gamma_rate |= 1 << pos;
        }
    }

    // The epsilon rate is composed by the least common digits of each position. This means that
    // this is just the complement of the gamma rate but we have to trim any extra left-bits to
    // have the right number of digits.
    //
    // This complement trick might fail if one of the positions only had one of the two
    // possible digits.
    let epsilon_rate = (!gamma_rate) & ((1 << (one_counts.len() as u32)) - 1);

    (gamma_rate, epsilon_rate)
}

fn main() -> Result<(), Error> {
    let input = std::fs::read_to_string("./input")?;
    // Iterator over the lines of the input file.
    let mut lines = input.lines();
    // We will store the input parsed as integers here.
    let mut numbers = Vec::with_capacity(lines.size_hint().0);

    // We need to special-case the first line to extract the number of digits.
    if let Some(first_line) = lines.next() {
        // The number of digits is just the length of the first line.
        let digits = first_line.len();

        // Panic if the input numbers are too long to fit into `u64` integers.
        assert!(
            digits <= u64::BITS.try_into().unwrap(),
            "numbers are too long"
        );

        // Parse the first line and push it to the numbers buffer if successful.
        numbers.push(u64::from_str_radix(first_line, 2)?);
        // Parse the remaining lines and push it to the numbers buffer if successful.
        for line in lines {
            numbers.push(u64::from_str_radix(line, 2)?);
        }
        // Count the number of one-digits in each position.
        let one_counts = count_ones(&numbers, digits);
        // Compute the rates.
        let (gamma_rate, epsilon_rate) = compute_rates(&one_counts, numbers.len());

        println!("Part 1: {}", gamma_rate * epsilon_rate);
    } else {
        eprintln!("Input is empty");
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test1() {
        let input = [
            0b00100, 0b11110, 0b10110, 0b10111, 0b10101, 0b01111, 0b00111, 0b11100, 0b10000,
            0b11001, 0b00010, 0b01010,
        ];

        let one_counts = count_ones(&input, 5);
        let (gamma_rate, epsilon_rate) = compute_rates(&one_counts, input.len());

        assert_eq!(0b10110, gamma_rate);
        assert_eq!(0b01001, epsilon_rate);
    }
}
