use std::collections::BTreeMap;

use Avichucho::*;

#[derive(Clone, Copy, Debug, Hash, PartialEq, Eq)]
enum Avichucho {
    A,
    B,
    C,
    D,
}

enum Tile {
    Wall,
    Hall,
    Door,
    Room(Avichucho),
}

struct State {
    tiles: BTreeMap<(usize, usize), Tile>,
    avichuchos: BTreeMap<(usize, usize), Avichucho>,
}

impl std::fmt::Display for State {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        for x in 1..=5 {
            for y in 1..=13 {
                let pos = (x, y);
                let chr = if let Some(tile) = self.tiles.get(&pos) {
                    match tile {
                        Tile::Wall => '#',
                        Tile::Door => '.',
                        Tile::Hall | Tile::Room(_) => {
                            if let Some(avichucho) = self.avichuchos.get(&pos) {
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

impl State {
    fn new(init: [Avichucho; 8]) -> Self {
        let mut tiles = BTreeMap::new();
        let mut avichuchos = BTreeMap::new();

        // Room A:
        avichuchos.insert((3, 4), init[0]);
        avichuchos.insert((4, 4), init[1]);
        // Room B:
        avichuchos.insert((3, 6), init[2]);
        avichuchos.insert((4, 6), init[3]);
        // Room C:
        avichuchos.insert((3, 8), init[4]);
        avichuchos.insert((4, 8), init[5]);
        // Room D:
        avichuchos.insert((3, 10), init[6]);
        avichuchos.insert((4, 10), init[7]);

        // Walls
        for y in 1..=13 {
            tiles.insert((1, y), Tile::Wall);
        }

        for y in 3..=11 {
            tiles.insert((5, y), Tile::Wall);
        }

        for y in 1..=3 {
            tiles.insert((3, y), Tile::Wall);
        }

        for y in 11..=13 {
            tiles.insert((3, y), Tile::Wall);
        }

        for x in 3..=4 {
            tiles.insert((x, 5), Tile::Wall);
            tiles.insert((x, 7), Tile::Wall);
            tiles.insert((x, 9), Tile::Wall);
        }

        tiles.insert((2, 1), Tile::Wall);
        tiles.insert((4, 3), Tile::Wall);
        tiles.insert((4, 11), Tile::Wall);
        tiles.insert((2, 13), Tile::Wall);

        // Hall & Doors
        for y in 2..=12 {
            if [4, 6, 8, 10].contains(&y) {
                tiles.insert((2, y), Tile::Door);
            } else {
                tiles.insert((2, y), Tile::Hall);
            }
        }

        // Rooms
        for (y, room) in [(4, A), (6, B), (8, C), (10, D)] {
            for x in 3..=4 {
                tiles.insert((x, y), Tile::Room(room));
            }
        }

        Self { tiles, avichuchos }
    }
}

fn main() {
    let state = State::new([B, A, C, D, B, C, D, A]);
    println!("{}", state);
}
