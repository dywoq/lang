namespace DywoqLang.Tools;

/// <summary>
/// Represents the AST node, or token
/// position with line, column and file.
/// </summary>
public record Position(int Line, int Column, int PPosition, string File);
