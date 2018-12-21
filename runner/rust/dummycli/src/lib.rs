extern "C" {
    fn __life_log(msg: *const u8, len: usize);
    fn __life_write_hchan(msg: *const u8, len: usize);
    fn __life_read_hchan(msg: *const u8, len: usize) -> usize;
}

use std::str;
use std::ptr;

fn life_read_chan() {
    let buf = [0; 1024];
    unsafe {
        let sz = __life_read_hchan(buf.as_ptr(), buf.len());
        let mut dst = Vec::with_capacity(sz);
        dst.set_len(sz);
        ptr::copy(buf.as_ptr(), dst.as_mut_ptr(), sz);
        life_log(&format!("sz: {}", sz));
        life_log(str::from_utf8(&dst).unwrap());
    };
}

fn life_write_chan(msg: &str) {
    unsafe {
        let bytes = msg.as_bytes();
        __life_write_hchan(bytes.as_ptr(), bytes.len())
    }
}

fn life_log(msg: &str) {
    unsafe {
        __life_log(msg.as_bytes().as_ptr(), msg.len());
    }
}

#[no_mangle]
pub extern "C" fn fib(n: i32) -> i32 {
	if n == 1 || n == 2 {
		1
	} else {
		fib(n - 1) + fib(n - 2)
	}
}

#[no_mangle]
pub extern "C" fn app_main() -> i32 {
    life_log("running...");
    life_write_chan("from insecure rust running inside a sandbox");
    life_read_chan();
	let result = fib(35);
    life_write_chan("from insecure rust running inside a sandbox (after fib(35))");
    result
}