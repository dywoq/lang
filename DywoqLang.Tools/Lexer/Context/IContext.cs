namespace DywoqLang.Tools.Lexer.Context;

public interface IContext
{
	/// <summary>
	/// Advances to the next position by n.
	/// </summary>
	public void Advance(int n);

	/// <summary>
	/// Returns the current character being process,
	/// or null.
	/// </summary>
	public char? Current();

	/// <summary>
	/// Returns the sub-string of the given input,
	/// or null.
	/// </summary>
	public string? Substring(int start, int end);

	/// <summary>
	/// Reports whether the lexer met the end of file.
	/// </summary>
	public bool EndOfFile();

	/// <summary>
	/// Returns the current position.
	/// </summary>
	public Position Position();

	/// <summary>
	/// Returns the future character, or null.
	public char? Peek();
}
