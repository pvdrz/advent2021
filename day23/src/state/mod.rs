mod board;
mod display;
mod iter;

use std::collections::BTreeMap;
use Amphipod::*;

pub use self::board::Board;

#[derive(Clone, Copy, Debug, Hash, PartialEq, Eq)]
pub enum Amphipod {
    A,
    B,
    C,
    D,
}

#[derive(Clone, Copy, Debug, Hash, PartialEq, Eq)]
pub enum Tile {
    Wall,
    Hall,
    Door,
    Room(Amphipod),
}

#[derive(Clone, Copy, Debug, Hash, PartialEq, Eq)]
enum Move {
    Hall,
    Room,
}

type Cost = usize;

impl Amphipod {
    pub fn step_cost(&self) -> Cost {
        match self {
            A => 1,
            B => 10,
            C => 100,
            D => 1000,
        }
    }

    pub fn cost(self, steps: usize) -> Cost {
        steps * self.step_cost()
    }
}

#[derive(Debug)]
pub struct State<'a> {
    amphipods: BTreeMap<(usize, usize), (Amphipod, Option<Move>)>,
    board: &'a Board,
}

impl<'a> State<'a> {
    pub fn new(init: [Amphipod; 8], board: &'a Board) -> Self {
        let mut amphipods = BTreeMap::new();

        // Room A:
        amphipods.insert((3, 4), (init[0], Some(Move::Hall)));
        amphipods.insert((4, 4), (init[1], Some(Move::Hall)));
        // Room B:
        amphipods.insert((3, 6), (init[2], Some(Move::Hall)));
        amphipods.insert((4, 6), (init[3], Some(Move::Hall)));
        // Room C:
        amphipods.insert((3, 8), (init[4], Some(Move::Hall)));
        amphipods.insert((4, 8), (init[5], Some(Move::Hall)));
        // Room D:
        amphipods.insert((3, 10), (init[6], Some(Move::Hall)));
        amphipods.insert((4, 10), (init[7], Some(Move::Hall)));

        Self { amphipods, board }
    }
}
