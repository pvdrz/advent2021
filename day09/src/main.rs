mod grid;

use anyhow::{Context, Error};

use crate::grid::Grid;

fn main() -> Result<(), Error> {
    let input = std::fs::read_to_string("./input").context("could not read input file.")?;

    let grid: Grid = input.parse().context("could not parse grid.")?;

    println!(
        "Part 1: {}",
        grid.get_mins()
            .into_iter()
            .map(|(i, j)| grid.get(i, j).unwrap() + 1)
            .sum::<u32>()
    );

    Ok(())
}

#[cfg(test)]
mod tests {
    use crate::grid::Grid;

    const EXAMPLE: &str = "2199943210
3987894921
9856789892
8767896789
9899965678";

    #[test]
    fn example_1() {
        let grid: Grid = EXAMPLE.parse().unwrap();

        assert_eq!(10, grid.n_cols(), "invalid number of columns");
        assert_eq!(5, grid.n_rows(), "invalid number of rows");

        let mins = grid.get_mins();
        assert_eq!(4, mins.len(), "invalid number of mins");

        let risk = mins
            .iter()
            .copied()
            .map(|(i, j)| grid.get(i, j).unwrap() + 1)
            .sum::<u32>();
        assert_eq!(15, risk, "invalid risk");

        for min in mins {
            let mut queue = vec![min];

            while let Some((i, j)) = queue.pop() {
                let value = grid.get(i, j).unwrap();
                if value > 9 {
                    continue;
                }
            }
        }
    }
}
