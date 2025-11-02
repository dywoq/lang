using DywoqLang.Tools.Lexer.Context;

namespace DywoqLang.Tools.Lexer.Token;

public interface IToken
{
	/// <summary>
	/// Returns the kind of the token.
	/// </summary>
	/// <returns></returns>
	public string Kind();

	/// <summary>
	/// Returns the tokenized character. May return false, 
	/// indicating this symbol doesn't meet the token's requirements,
	/// then, the lexer should try other token.
	/// </summary>
	public (TokenData, bool) Tokenize(IContext context);
}
