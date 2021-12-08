use std::ops::RangeInclusive;

use anyhow::{Context, Error};

fn parse_input(input: &str) -> Result<Vec<Segment>, Error> {
    let mut segments = Vec::new();

    for line in input.trim().lines() {
        if let Some(segment) = Segment::parse(line)? {
            segments.push(segment);
        }
    }

    Ok(segments)
}

#[derive(Debug, Clone, Copy)]
enum SegmentKind {
    Horizontal,
    Vertical,
}

#[derive(Debug)]
struct Segment {
    interval: RangeInclusive<usize>,
    intersect: usize,
    kind: SegmentKind,
}

impl Segment {
    fn parse(input: &str) -> Result<Option<Self>, Error> {
        let (fst, snd) = input.trim().split_once(" -> ").context("missing ` -> `")?;

        let (x1, y1) = fst
            .split_once(',')
            .context("missing `,` for the first pair of points")?;

        let (x2, y2) = snd
            .split_once(',')
            .context("missing `,` for the second pair of points")?;

        Ok(Self::new(
            (x1.parse()?, y1.parse()?),
            (x2.parse()?, y2.parse()?),
        ))
    }

    /// Creates a new segment from two points.
    fn new((x1, y1): (usize, usize), (x2, y2): (usize, usize)) -> Option<Self> {
        // If both `x` coordinates are equal, the segment is vertical.
        if x1 == x2 {
            Some(Self {
                interval: if y1 < y2 { y1..=y2 } else { y2..=y1 },
                intersect: x1,
                kind: SegmentKind::Vertical,
            })
        // If both `y` coordinates are equal, the segment is horizontal.
        } else if y1 == y2 {
            Some(Self {
                interval: if x1 < x2 { x1..=x2 } else { x2..=x1 },
                intersect: y1,
                kind: SegmentKind::Horizontal,
            })
        } else {
            None
        }
    }

    /// Computes the overlap between the intervals of the segments.
    fn interval_overlaps<'a>(&'a self, other: &'a Self) -> impl Iterator<Item = usize> + 'a {
        self.interval
            .clone()
            .filter(|&z| other.interval_contains(z))
    }

    /// Checks if the interval of the segment contains one point.
    fn interval_contains(&self, z: usize) -> bool {
        self.interval.contains(&z)
    }

    /// Computes the intersection of two segments.
    fn intersection(&self, other: &Self) -> Vec<(usize, usize)> {
        // Check the cases where the intersection could be non-empty
        match (self.kind, other.kind) {
            // If both segments are horizontal and their intersects are equal
            (SegmentKind::Horizontal, SegmentKind::Horizontal)
                if self.intersect == other.intersect =>
            {
                // Compute the overlap between the intervals and use the intersect as the `y`
                // coordinate.
                self.interval_overlaps(other)
                    .map(|z| (z, self.intersect))
                    .collect()
            }
            // If both segments are vertical and their intersects are equal
            (SegmentKind::Vertical, SegmentKind::Vertical)
                if (self.intersect == other.intersect) =>
            {
                // Compute the overlap between the intervals and use the intersect as the `x`
                // coordinate.
                self.interval_overlaps(other)
                    .map(|z| (self.intersect, z))
                    .collect()
            }
            // If the first segment is vertical, the second one is horizontal, the interval of the
            // first contains the intersect of the second and viceversa.
            (SegmentKind::Vertical, SegmentKind::Horizontal)
                if self.interval_contains(other.intersect)
                    && other.interval_contains(self.intersect) =>
            {
                // The first segment is vertical, meaning that `x` is constant.
                vec![(self.intersect, other.intersect)]
            }
            // The converse of the previous case.
            (SegmentKind::Horizontal, SegmentKind::Vertical)
                if self.interval_contains(other.intersect)
                    && other.interval_contains(self.intersect) =>
            {
                // The second segment is vertical, meaning that `x` is constant.
                vec![(other.intersect, self.intersect)]
            }
            // The segments do not overlap.
            _ => Vec::new(),
        }
    }
}

fn count_intersects(segments: &[Segment]) -> usize {
    let mut points = Vec::new();

    for (i, seg1) in segments.iter().enumerate() {
        if let Some(segments) = segments.get((i + 1)..) {
            for seg2 in segments {
                for point in seg1.intersection(seg2) {
                    if let Err(i) = points.binary_search(&point) {
                        points.insert(i, point);
                    }
                }
            }
        }
    }

    points.len()
}

fn main() -> Result<(), Error> {
    let segments = parse_input(&std::fs::read_to_string("./input")?)?;

    println!("Part 1: {}", count_intersects(&segments));

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test1() {
        let segments = parse_input(
            "0,9 -> 5,9
    8,0 -> 0,8
    9,4 -> 3,4
    2,2 -> 2,1
    7,0 -> 7,4
    6,4 -> 2,0
    0,9 -> 2,9
    3,4 -> 1,4
    0,0 -> 8,8
    5,5 -> 8,2",
        )
        .unwrap();

        assert_eq!(5, count_intersects(&segments));
    }
}
