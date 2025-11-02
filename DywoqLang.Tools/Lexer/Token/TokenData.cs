namespace DywoqLang.Tools.Lexer.Token;

/// <summary>
/// Represents the token data with literal, kind and position.
/// </summary>
public record TokenData(string Literal, string Kind, Position Position);
