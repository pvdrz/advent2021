pub(crate) const fn to_binary(c: u8) -> [bool; 4] {
    match c.to_ascii_uppercase() {
        b'0' => [false, false, false, false],
        b'1' => [false, false, false, true],
        b'2' => [false, false, true, false],
        b'3' => [false, false, true, true],
        b'4' => [false, true, false, false],
        b'5' => [false, true, false, true],
        b'6' => [false, true, true, false],
        b'7' => [false, true, true, true],
        b'8' => [true, false, false, false],
        b'9' => [true, false, false, true],
        b'A' => [true, false, true, false],
        b'B' => [true, false, true, true],
        b'C' => [true, true, false, false],
        b'D' => [true, true, false, true],
        b'E' => [true, true, true, false],
        b'F' => [true, true, true, true],
        _ => panic!("char is not a hexadecimal digit"),
    }
}
