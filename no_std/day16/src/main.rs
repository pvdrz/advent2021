#![no_std]
#![no_main]
mod counting_iter;
mod decode;
mod hexa;

use libc_print::std_name::println;

pub type Error = &'static str;

const INPUT: &str ="3600888023024c01150044c0118330a440118330e44011833085c0118522008c29870";

#[no_mangle]
pub extern "C" fn main(_argc: isize, _argv: *const *const u8) -> isize {
    println!("INPUT: {}", INPUT);
    let (version_sum, result) = decode::evall(INPUT).unwrap();
    println!("{} {}", version_sum, result);

    0
}

#[panic_handler]
fn my_panic(info: &core::panic::PanicInfo) -> ! {
    println!("oh fuck!");
    println!("{:?}", info);
    loop {}
}

