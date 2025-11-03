namespace DywoqLang.Tools.Lexer.Token;

/// <summary>
/// An exception class is thrown by the tokens and their tokenizers
/// if they got something they didn't expect.
/// </summary>
public class ExpectedException : Exception
{
	public ExpectedException() { }
	public ExpectedException(string message) : base(message) { }
	public ExpectedException(string message, Exception inner) : base(message, inner) { }


	/// <summary>
	/// Returns the ExpectedException with the automatically formatted message
	/// with no position.
	/// </summary>
	public static ExpectedException From(string message) =>
			new($"ExpectedException: {message}; source is unknown");

	/// <summary>
	/// Returns the ExpectedException with the automatically formatted message
	/// with position.
	/// </summary>
	public static ExpectedException From(string message, Position position) =>
		new($"ExpectedException: {message}; source: is {position.File}:{position.Line}:{position.Column}");
}
