using DywoqLang.Tools.Lexer.Context;

namespace DywoqLang.Tools.Lexer.Token;

/// <summary>
/// Presents the integer, or float token in DywoqLang.
/// Its presentation: 8723 (integer), 3.14 (float).
/// </summary>
/// <see cref="IToken"/>
public class Number : IToken
{
	public (TokenData?, bool) Tokenize(IContext context)
	{
		var current = context.Current() ?? throw new NullReferenceException();
		if (!char.IsDigit(current))
			return (null, false);

		context.Advance(1);

		var start = context.Position().PPosition;
		while (!context.EndOfFile())
		{
			current = context.Current() ?? throw new NullReferenceException();
			if (!char.IsDigit(current))
			{
				break;
			}
			context.Advance(1);
		}

		if (context.Peek() == null || context.Peek() != '.')
		{
			// we place here directly to prevent CS0136: https://learn.microsoft.com/en-us/dotnet/csharp/misc/cs0136?f1url=%3FappId%3Droslyn%26k%3Dk(CS0136)
			return (new(context.Substring(start, context.Position().PPosition) ?? throw new NullReferenceException(), KindGetter.AsInteger(), context.Position()), true);
		}

		context.Advance(1);

		// parse numbers
		current = context.Current() ?? throw new NullReferenceException();
		if (context.Current() == null || !char.IsDigit(current))
		{
			throw ExpectedException.From("expected digit after float", context.Position());
		}

		while (!context.EndOfFile())
		{
			current = context.Current() ?? throw new NullReferenceException();
			if (!char.IsDigit(current))
			{
				break;
			}
			context.Advance(1);
		}

		var literal = context.Substring(start, context.Position().PPosition) ?? throw new NullReferenceException();
		return (new(literal, KindGetter.AsFloat(), context.Position()), true);
	}
}
