use super::{Cost, Move, State, Tile};

impl<'a> State<'a> {
    pub fn expand(self) -> Vec<(Self, Cost)> {
        let State { amphipods, board } = self;
        let mut expanded_states = Vec::new();
        for ((x, y), (amphipod, next_move)) in amphipods.clone() {
            match next_move {
                Some(Move::Hall) => {
                    for (pos, tile_ty) in board.iter() {
                        if *tile_ty == Tile::Hall {
                            // A candidate for the next move. For that it needs to be reachable.
                            let (tile_x, tile_y) = *pos;

                            let is_reachable = (tile_x..x)
                                .all(|x| !amphipods.contains_key(dbg!(&(x, y))))
                                && (tile_y.min(y)..=tile_y.max(y))
                                    .all(|y| !amphipods.contains_key(dbg!(&(tile_x, y))));

                            if is_reachable {
                                let mut new_amphipods = amphipods.clone();
                                new_amphipods.remove(&(x, y));
                                let cost =
                                    amphipod.cost(tile_y.max(y) - tile_y.min(y) + x - tile_x);
                                new_amphipods.insert(*pos, (amphipod, Some(Move::Room)));
                                expanded_states.push((
                                    State {
                                        amphipods: new_amphipods,
                                        board: board.clone(),
                                    },
                                    cost,
                                ));
                            }
                        }
                    }
                }
                Some(Move::Room) => {
                    for (pos, tile_ty) in board.iter() {
                        if *tile_ty == Tile::Room(amphipod) {
                            // A candidate for the next move. For that it needs to be reachable.
                            let (tile_x, tile_y) = *pos;

                            let (max_y, min_y) = if tile_y > y {
                                (tile_y, y + 1)
                            } else {
                                (y - 1, tile_y)
                            };
                            let is_reachable = (x + 1..=tile_x)
                                .all(|x| !amphipods.contains_key(dbg!(&(x, tile_y))))
                                && (min_y..=max_y).all(|y| !amphipods.contains_key(dbg!(&(x, y))));

                            if is_reachable {
                                let mut new_amphipods = amphipods.clone();
                                new_amphipods.remove(&(x, y));
                                let cost = amphipod.cost(max_y - min_y + tile_x - x);
                                new_amphipods.insert(*pos, (amphipod, None));
                                expanded_states.push((
                                    State {
                                        amphipods: new_amphipods,
                                        board: board.clone(),
                                    },
                                    cost,
                                ));
                            }
                        }
                    }
                }
                None => todo!(),
            }
        }

        expanded_states
    }
}
