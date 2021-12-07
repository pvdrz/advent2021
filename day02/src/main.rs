use anyhow::{Context, Error};
use std::{fs::read_to_string, num::ParseIntError, str::FromStr};

#[derive(Debug)]
enum Command {
    Forward(usize),
    Down(usize),
    Up(usize),
}

#[derive(Debug)]
enum CommandParseError {
    MissingSpace,
    UnknownCommand,
    Int(ParseIntError),
}

impl From<ParseIntError> for CommandParseError {
    fn from(err: ParseIntError) -> Self {
        Self::Int(err)
    }
}

impl std::fmt::Display for CommandParseError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            Self::MissingSpace => write!(f, "missing space separator"),
            Self::UnknownCommand => write!(f, "unknown command prefix"),
            Self::Int(err) => err.fmt(f),
        }
    }
}

impl std::error::Error for CommandParseError {}

impl FromStr for Command {
    type Err = CommandParseError;

    fn from_str(s: &str) -> Result<Self, Self::Err> {
        let (head, tail) = s
            .split_once(' ')
            .ok_or_else(|| CommandParseError::MissingSpace)?;
        let value = usize::from_str(tail)?;

        let command = match head {
            "forward" => Self::Forward(value),
            "down" => Self::Down(value),
            "up" => Self::Up(value),
            _ => return Err(CommandParseError::UnknownCommand),
        };

        Ok(command)
    }
}

fn final_position(commands: &[Command]) -> (usize, usize) {
    let (mut x, mut y) = (0, 0);

    for command in commands {
        match command {
            Command::Forward(dx) => x += dx,
            Command::Down(dy) => y += dy,
            Command::Up(dy) => y -= dy,
        }
    }

    (x, y)
}

fn final_position_with_aim(commands: &[Command]) -> (usize, usize, usize) {
    let (mut x, mut y, mut aim) = (0, 0, 0);

    for command in commands {
        match command {
            Command::Forward(dx) => {
                x += dx;
                y += aim * dx;
            }
            Command::Down(daim) => aim += daim,
            Command::Up(daim) => aim -= daim,
        }
    }

    (x, y, aim)
}

fn main() -> Result<(), Error> {
    let input = read_to_string("./input")
        .context("could not read input file")?
        .lines()
        .map(|line| line.parse().context("could not parse input line"))
        .collect::<Result<Vec<Command>, _>>()?;

    let (x, y) = final_position(&input);
    println!("Part 1: {}", x * y);

    let (x, y, _aim) = final_position_with_aim(&input);
    println!("Part 2: {}", x * y);

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test1() {
        let input = [
            "forward 5",
            "down 5",
            "forward 8",
            "up 3",
            "down 8",
            "forward 2",
        ]
        .into_iter()
        .map(|s| s.parse().unwrap())
        .collect::<Vec<Command>>();

        assert_eq!((15, 10), final_position(&input));
    }

    #[test]
    fn test2() {
        let input = [
            "forward 5",
            "down 5",
            "forward 8",
            "up 3",
            "down 8",
            "forward 2",
        ]
        .into_iter()
        .map(|s| s.parse().unwrap())
        .collect::<Vec<Command>>();

        let (x, y, _aim) = final_position_with_aim(&input);
        assert_eq!((15, 60), (x, y));
    }
}
