use std::str::FromStr;

use anyhow::{Context, Error};

/// The number of rows in each board.
const ROWS: usize = 5;
/// The number of columns in each board.
const COLS: usize = 5;

/// A bingo board.
#[derive(Debug)]
struct Board {
    inner: [[(usize, bool); COLS]; ROWS],
}

impl Board {
    /// Mark a number in the board and return `true` if the board won or `false` otherwise.
    fn mark_number(&mut self, number: usize) -> bool {
        // Here we will store the position of the cells that were marked so they can be checked
        // later.
        let mut marked_cells = Vec::new();

        for (i, row) in self.inner.iter_mut().enumerate() {
            for (j, (cell, marked)) in row.iter_mut().enumerate() {
                // If the current cell has the number and it has not been marked yet, mark it and
                // push it to `marked_cells`.
                if *cell == number && !*marked {
                    *marked = true;
                    marked_cells.push((i, j));
                    break;
                }
            }
        }

        // Check if any of the marked cells caused this board to win
        marked_cells
            .into_iter()
            .any(|(i, j)| self.check_if_won(i, j))
    }

    /// Check if the current board won because of the `i`th row or the `j`th column.
    fn check_if_won(&self, i: usize, j: usize) -> bool {
        self.inner[i].iter().all(|(_, marked)| *marked)
            || self
                .inner
                .iter()
                .map(|row| row[j])
                .all(|(_, marked)| marked)
    }

    /// Compute the sum of the unmarked cells.
    fn unmarked_sum(&self) -> usize {
        let mut count = 0;
        for row in &self.inner {
            for (cell, marked) in row {
                if !marked {
                    count += *cell;
                }
            }
        }
        count
    }
}

impl FromStr for Board {
    type Err = Error;

    fn from_str(input: &str) -> Result<Self, Self::Err> {
        let mut inner = [[(0usize, false); COLS]; ROWS];
        let mut lines = input.lines();

        for row in &mut inner {
            // Split each line of the input by whitespaces and parse each chunk as an integer.
            let mut nums = lines.next().context("missing row")?.split_whitespace();
            for (cell, _) in row.iter_mut() {
                *cell = nums.next().context("missing number")?.parse()?;
            }
        }

        Ok(Self { inner })
    }
}

fn main() -> Result<(), Error> {
    let input = std::fs::read_to_string("./input").context("could not read input file")?;

    // `head` has the numbers to be drawn and `tail` the boards.
    let (head, tail) = input.split_once("\n\n").context("invalid header")?;

    // Parse the numbers splitting `head` by commas and parsing each chunk as an integer.
    let numbers = head
        .split(',')
        .map(|s| s.parse())
        .collect::<Result<Vec<usize>, _>>()
        .context("could not parse numbers to be drawn")?;

    // Parse the boards splitting `tail` every two new lines.
    let mut boards = tail
        .split("\n\n")
        .map(|s| s.parse())
        .collect::<Result<Vec<Board>, _>>()
        .context("could not parse boards")?;

    for number in numbers {
        // Here we store the indices of the boards that must be deleted because they already won.
        let mut indices_to_delete = Vec::new();

        for (index, board) in boards.iter_mut().enumerate() {
            // If the board won while drawing this number, we print its score and then add it to
            // the deletion buffer.
            if board.mark_number(number) {
                println!("Board won: {}", board.unmarked_sum() * number);
                // We add it in order because deletion must be done starting with the largest index
                // first.
                if let Err(i) = indices_to_delete.binary_search(&index) {
                    indices_to_delete.insert(i, index);
                }
            }
        }

        // Delete the winning boards starting from the largest index to avoid index shifting.
        for index in indices_to_delete.into_iter().rev() {
            boards.remove(index);
        }
    }

    Ok(())
}
