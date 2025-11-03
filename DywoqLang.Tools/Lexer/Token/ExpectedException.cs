namespace DywoqLang.Tools.Lexer.Token;

/// <summary>
/// An exception class is thrown by the tokens and their tokenizers
/// if they got something they didn't expect.
/// </summary>
public class ExpectedException : Exception
{
	public Position? Position { get; set; }

	public ExpectedException() { }
	public ExpectedException(string message) : base(message) { }
	public ExpectedException(string message, Position position) : base(message)
	{
		Position = position;
	}

	public ExpectedException(string message, Position position, Exception inner) : base(message, inner)
	{
		Position = position;
	}

	/// <summary>
	/// Returns the ExpectedException with the automatically formatted message
	/// with no position.
	/// </summary>
	public static ExpectedException From(string message) =>
			new($"{message}; source is unknown");

	/// <summary>
	/// Returns the ExpectedException with the automatically formatted message
	/// with position.
	/// </summary>
	public static ExpectedException From(string message, Position position)
		=> new($"{message}; source: is {position.File}:{position.Line}:{position.Column}", position);
}
