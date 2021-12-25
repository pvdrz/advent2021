use std::collections::{hash_map::Entry, BinaryHeap, HashMap};

use super::{Cost, Move, State, Tile};

impl<'a> State<'a> {
    // Every reachable move from the current state.
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
                                .all(|x| !amphipods.contains_key(&(x, y)))
                                && (tile_y.min(y)..=tile_y.max(y))
                                    .all(|y| !amphipods.contains_key(&(tile_x, y)));

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
                                .all(|x| !amphipods.contains_key(&(x, tile_y)))
                                && (min_y..=max_y).all(|y| !amphipods.contains_key(&(x, y)));

                            if is_reachable {
                                let worth_doing = amphipods.contains_key(&(tile_x + 1, tile_y))
                                    || board
                                        .get(&(tile_x + 1, tile_y))
                                        .map(|t| *t == Tile::Wall)
                                        .unwrap();

                                if !worth_doing {
                                    continue;
                                }
                                let mut new_amphipods = amphipods.clone();
                                new_amphipods.remove(&(x, y));
                                let cost =
                                    amphipod.cost(tile_y.max(y) - tile_y.min(y) + tile_x - x);
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
                None => {}
            }
        }

        expanded_states
    }
}

#[derive(PartialEq, Eq)]
struct Queued<'a> {
    state: State<'a>,
    cost: Cost,
}

impl<'a> PartialOrd for Queued<'a> {
    fn partial_cmp(&self, other: &Self) -> Option<std::cmp::Ordering> {
        Some(self.cmp(other))
    }
}

impl<'a> Ord for Queued<'a> {
    fn cmp(&self, other: &Self) -> std::cmp::Ordering {
        other.cost.cmp(&self.cost)
    }
}

pub fn find_min(init: State) -> Cost {
    let mut best_costs = HashMap::from([(init.clone(), 0)]);
    let mut prev: HashMap<State, State> = HashMap::new();
    let mut queue = BinaryHeap::from([Queued {
        state: init,
        cost: 0,
    }]);

    while let Some(Queued { state, cost }) = queue.pop() {
        if state.is_final() {
            let mut path = vec![state];
            while let Some(prev) = prev.get(path.last().unwrap()) {
                path.push(prev.clone());
            }
            for state in path.into_iter().rev() {
                println!("\n{}", state);
            }
            return cost;
        }
        for (expanded, expansion_cost) in state.clone().expand() {
            let alt = cost + expansion_cost;
            match best_costs.entry(expanded.clone()) {
                Entry::Occupied(mut entry) => {
                    if alt < *entry.get() {
                        entry.insert(alt);
                        prev.insert(expanded.clone(), state.clone());
                        queue.push(Queued {
                            state: expanded,
                            cost: alt,
                        })
                    }
                }
                Entry::Vacant(entry) => {
                    entry.insert(alt);
                    prev.insert(expanded.clone(), state.clone());
                    queue.push(Queued {
                        state: expanded,
                        cost: alt,
                    })
                }
            }
        }
    }

    panic!("Path not found");
}
