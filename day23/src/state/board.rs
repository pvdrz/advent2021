use std::collections::BTreeMap;

use super::{Amphipod::*, Tile};

#[derive(Debug)]
pub struct Board {
    tiles: BTreeMap<(usize, usize), Tile>,
}

impl std::ops::Deref for Board {
    type Target = BTreeMap<(usize, usize), Tile>;

    fn deref(&self) -> &Self::Target {
        &self.tiles
    }
}

impl std::ops::DerefMut for Board {
    fn deref_mut(&mut self) -> &mut Self::Target {
        &mut self.tiles
    }
}

impl Default for Board {
    fn default() -> Self {
        let mut tiles = BTreeMap::new();

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

        Self { tiles }
    }
}
