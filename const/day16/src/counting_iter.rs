use crate::hexa;

pub(crate) struct CountingIter<'a> {
    next_calls: usize,
    head: [bool; 4],
    buffer: &'a [u8],
    pos: usize,
}

impl<'a> CountingIter<'a> {
    pub(crate) const fn new(buffer: &'a [u8]) -> Self {
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

    pub(crate) const fn calls(&self) -> usize {
        self.next_calls
    }

    pub(crate) const fn next(mut self) -> (Self, bool) {
        self.next_calls += 1;

        let item = self.head[self.pos];

        if self.pos == 3 {
            if self.buffer.len() != 0 {
                let ego = Self::new(self.buffer);
                self = Self {
                    next_calls: self.next_calls,
                    ..ego
                }
            }
        } else {
            self.pos += 1;
        }

        (self, item)
    }
}
