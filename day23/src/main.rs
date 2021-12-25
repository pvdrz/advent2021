use state::find_min;

use crate::state::Amphipod::*;
use crate::state::Board;
use crate::state::State;

mod state;
fn main() {
    // let board = Board::default();
    // let state = State::new([B, A, C, D, B, C, D, A], &board); // example

    // let state = State::new([D, B, D, A, C, B, C, A], &board); // divi
    // let state = State::new([C, B, A, A, D, B, D, C], &board); // chris

    // println!("{}", state);
    // let min_cost = find_min(state);
    // println!("Min cost: {}.", min_cost);

    let board = Board::unfolded();
    // let state = State::new_unfolded([B, A, C, D, B, C, D, A], &board);
    let state = State::new_unfolded([D, B, D, A, C, B, C, A], &board);
    // let state = State::new_unfolded([C, B, A, A, D, B, D, C], &board);
    println!("{}", state);
    let min_cost = find_min(state);
    println!("Min cost: {}.", min_cost);
}
