use crate::state::Amphipod::*;
use crate::state::Board;
use crate::state::State;

mod state;
fn main() {
    let board = Board::default();
    let state = State::new([B, A, C, D, B, C, D, A], &board);
    println!("{}", state);
    for (state, cost) in state.expand() {
        println!("Next state with cost {}\n{}", cost, state);
    }
}
