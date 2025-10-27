# Types

## Basic types
| Type   | Range                                                   | Additional                    |
| ------ | ------------------------------------------------------- | ----------------------------- |
| `i8`   | -128 to 127                                             | No additional information     |
| `i16`  | -32,768 to 32,767                                       | No additional information     |
| `i32`  | -2,147,483,648 to 2,147,483,647                         | No additional information     |
| `i64`  | -9,223,372,036,854,775,808 to 9,223,372,036,854,775,807 | No additional information     |
| `u8`   | 0 to 255                                                | No additional information     |
| `u16`  | 0 to 65,535                                             | No additional information     |
| `u32`  | 0 to 4,294,967,295                                      | No additional information     |
| `u64`  | 0 to 18,446,744,073,709,551,615                         | No additional information     |
| `str`  | Variable size (limited by available memory)             | UTF-8 encoded characters      |
| `bool` | 0 to 1                                                  | Constants are `true`, `false` |
| `void` | No range                                                | No additional information     |

## Types associated with memory
| Type   | Range                           | Additional     |
| ------ | ------------------------------- | -------------- |
| `uptr` | 0 to 18,446,744,073,709,551,615 | Alias of `u64` |

## Floating-point types
| Type  | Range                                                             | Additional                |
| ----- | ----------------------------------------------------------------- | ------------------------- |
| `f32` | ±1.401298464324817e-45 to ±3.4028234663852886e+38 (approximately) | No additional information |
| `f64` | ±4.9406564584124654e-324 to ±1.7976931348623157e+308              | No additional information |

## 128 bits integer
| Type   | Range               | Additional                |
| ------ | ------------------- | ------------------------- |
| `i128` | -2^127 to 2^127 - 1 | No additional information |
| `u128` | 0 to 2^128 - 1      | No additional information |

## Fixed-point types
| Type    | Range                                                    | Additional                        |
| ------- | -------------------------------------------------------- | --------------------------------- |
| `fix64` | -2^32 to 2^32 - 1 (with 32 bits for the fractional part) | Fixed-point 64-bit representation |
