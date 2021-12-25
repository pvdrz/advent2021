mod board;
mod display;

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

    pub fn expand(self, board: &Board) -> Vec<(Self, Cost)> {
        let State { board, amphipods } = self;
        let original_ocupied: Vec<_> = amphipods.keys().cloned().collect();
        let expanded_states = Vec::new();
        for ((x, y), (amphipod, next_move)) in amphipods {
            match next_move {
                Some(Move::Hall) => {
                    // this amphipod can move to a hall. Find all available halls and compute the
                    // cost of moving it there.
                    for ((tile_x, tile_y), tyle_ty) in &board {
                        if *tyle_ty == Tile::Hall {
                            // this distance is manhatan.
                            let steps =
                                x.max(*tile_x) - x.min(*tile_x) + y.max(*tile_y) - y.min(*tile_y);
                            let cost = amphipod.cost(steps);
                            let next_state = State {
                                board: board.clone(),
                                amphipods,
                            };
                        }
                    }
                }
                Some(Move::Room) => todo!(),
                None => {
                    // Can't move more, another one will need to move
                }
            }
        }
        expanded_states
    }
}
