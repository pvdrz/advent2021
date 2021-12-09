use anyhow::{anyhow, Context, Error};

#[derive(Debug)]
pub(crate) struct Grid {
    rows: Vec<Vec<u32>>,
    n_cols: usize,
}

impl Grid {
    #[inline]
    pub(crate) fn n_cols(&self) -> usize {
        self.n_cols
    }

    #[inline]
    pub(crate) fn n_rows(&self) -> usize {
        self.rows.len()
    }

    pub(crate) fn get(&self, i: usize, j: usize) -> Option<u32> {
        self.rows.get(j)?.get(i).copied()
    }
    pub(crate) fn contains(&self, i: usize, j: usize) -> bool {
        i < self.n_cols() && j < self.n_rows()
    }

    pub(crate) fn get_neighbors(&self, i: usize, j: usize) -> impl Iterator<Item = (usize, usize)> {
        let up = j
            .checked_sub(1)
            .map(|j| (i, j))
            .filter(|&(i, j)| self.contains(i, j));

        let left = i
            .checked_sub(1)
            .map(|i| (i, j))
            .filter(|&(i, j)| self.contains(i, j));

        let down = j
            .checked_add(1)
            .map(|j| (i, j))
            .filter(|&(i, j)| self.contains(i, j));

        let right = i
            .checked_add(1)
            .map(|i| (i, j))
            .filter(|&(i, j)| self.contains(i, j));

        up.into_iter().chain(left).chain(down).chain(right)
    }

    pub(crate) fn get_mins(&self) -> Vec<(usize, usize)> {
        let mut mins = Vec::new();

        for j in 0..self.n_rows() {
            for i in 0..self.n_cols() {
                let value = self.get(i, j).unwrap();

                if self
                    .get_neighbors(i, j)
                    .all(|(i, j)| value < (self.get(i, j).unwrap()))
                {
                    mins.push((i, j));
                }
            }
        }

        mins
    }
}

impl std::str::FromStr for Grid {
    type Err = Error;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let mut lines = s.lines().enumerate();

        let (_, first_line) = lines.next().context("input has no lines.")?;

        let first_row = first_line
            .char_indices()
            .map(|(col_idx, c)| {
                c.to_digit(10)
                    .with_context(|| format!("could not parse integer at (0, {}).", col_idx))
            })
            .collect::<Result<Vec<_>, _>>()?;

        let n_cols = first_row.len();

        let mut rows = vec![first_row];

        for (row_idx, line) in lines {
            let row = line
                .char_indices()
                .map(|(col_idx, c)| {
                    c.to_digit(10).with_context(|| {
                        format!("could not parse integer at ({}, {}).", row_idx, col_idx)
                    })
                })
                .collect::<Result<Vec<_>, _>>()?;

            if row.len() != n_cols {
                return Err(anyhow!(
                    "the row at line {} has {} elements, {} were expected.",
                    row_idx,
                    row.len(),
                    n_cols
                ));
            }

            rows.push(row);
        }

        Ok(Self { rows, n_cols })
    }
}
