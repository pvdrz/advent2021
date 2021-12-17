use crate::hexa;

pub struct CountingIter<'a> {
    next_calls: usize,
    head: [bool; 4],
    buffer: &'a [u8],
    pos: usize,
}

impl<'a> CountingIter<'a> {
    pub const fn new(buffer: &'a [u8]) -> Self {
        match buffer.split_first() {
            Some((&head, buffer)) => CountingIter {
                next_calls: 0,
                buffer,
                head: hexa::to_binary(head),
                pos: 0,
            },
            None => panic!("Buffer is empty"),
        }
    }

    pub const fn len(&self) -> usize {
        4 * self.buffer.len() + (4 - self.pos)
    }

    pub const fn calls(&self) -> usize {
        self.next_calls
    }

    pub const fn next(mut self, msg: &'static str) -> (Self, bool) {
        self.next_calls += 1;

        let item = self.head[self.pos];

        if self.pos == 3 {
            if self.buffer.len() != 0 {
                let ego = Self::new(self.buffer);
                self = Self {
                    next_calls: self.next_calls,
                    ..ego
                }
            } else {
                panic!("{}", msg);
            }
        } else {
            self.pos += 1;
        }

        (self, item)
    }
}
