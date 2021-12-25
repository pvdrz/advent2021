use super::{Amphipod::*, State, Tile};

impl<'a> std::fmt::Display for State<'a> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        write!(f, "x ")?;
        for y in 1..=13 {
            write!(f, "{}", y % 10)?;
        }
        write!(f, " <- y \n")?;
        for x in 1..=5 {
            write!(f, "{} ", x)?;
            for y in 1..=13 {
                let pos = (x, y);
                let chr = if let Some(tile) = self.board.get(&pos) {
                    match tile {
                        Tile::Wall => '█',
                        Tile::Door => '.',
                        Tile::Hall | Tile::Room(_) => {
                            if let Some((avichucho, _)) = self.amphipods.get(&pos) {
                                match avichucho {
                                    A => 'A',
                                    B => 'B',
                                    C => 'C',
                                    D => 'D',
                                }
                            } else {
                                ' '
                            }
                        }
                    }
                } else {
                    ' '
                };
                chr.fmt(f)?;
            }
            '\n'.fmt(f)?;
        }

        Ok(())
    }
}
