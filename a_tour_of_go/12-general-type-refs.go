/* 

1. General and Types
These work for almost any type of variable.
• %T (Type): Prints the exact type of the variable in Go (e.g., int, string, main.User). It is very useful for debugging code where the type is not obvious.
• %% (Literal percent sign): Prints the % symbol without expecting a variable.
• %t (Booleans): Prints the word true or false (has the same effect as %v for booleans).

2. Integers
Used to format whole numbers in different mathematical bases or representations.
• %d (Decimal): Standard base 10 (e.g., 42).
• %b (Binary): Base 2 (e.g., 101010).
• %o / %O (Octal): Base 8. %O adds the 0o prefix (e.g., 0o52).
• %x / %X (Hexadecimal): Base 16. %x uses lowercase letters (2a) and %X uses uppercase (2A).
• %c (Character): Displays the Unicode character corresponding to the number, without quotes (e.g., 65 becomes A).
• %U (Unicode format): Displays in the official Unicode standard format (e.g., U+0041).

3. Floating-Point and Complex Numbers
Used for float32, float64, and complex numbers.
• %f / %F (Standard decimal): Decimal point without an exponent (e.g., 123.456000). You can control the decimal places using %.2f (e.g., 123.46).
• %e / %E (Scientific notation): Displays with an exponent (e.g., 1.234560e+02).
• %g / %G (Compact): Automatically chooses between %f and %e to keep the output shorter and more readable.
• %b (Binary exponent): Shows the scientific representation in base 2.

4. Strings and Byte Slices ([]byte)
• %s (String): Prints the raw text, without quotes (e.g., Hello).
• %x / %X (Hexadecimal): Converts each byte of the string into two hexadecimal digits (e.g., 48656c6c6f). Widely used to display hashes or encrypted passwords.

5. Pointers and Memory
• %p (Pointer): Shows the memory address where the variable is stored, in hexadecimal format with the 0x prefix (e.g., 0xc0000140a8).
Extra Modifiers (Width and Alignment)
You can combine numbers and symbols between the % and the letter to align text in tables:
• %6d: Reserves 6 spaces and right-aligns the number.
• %-6d: Reserves 6 spaces and left-aligns the number.
• %06d: Pads empty spaces on the left with zeros (e.g., 42 becomes 000042).
*/