package goluaparser

func (this *LuaParser) Unop() error {
	/*@bgen(jjtree) Unop */
	jjtn000 := NewSimpleNode(JJTUNOP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	var err error
	cv := this.jj_ntk
	if cv == -1 {
		cv, err = this.jj_ntk_f()
		if err != nil {
			return err
		}
	}
	switch cv {
	case 83:
		{
			_, err = this.jj_consume_token(83)
			break
		}
	case NOT:
		{
			_, err = this.jj_consume_token(NOT)
			break
		}
	case 69:
		{
			_, err = this.jj_consume_token(69)
			break
		}
	default:
		this.jj_consume_token(-1)
		err = newParseException("unop error")
	}
	return err
}

func (this *LuaParser) Binop() error { /*@bgen(jjtree) Binop */
	jjtn000 := NewSimpleNode(JJTBINOP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	var err error
	cv := this.jj_ntk
	if cv == -1 {
		cv, err = this.jj_ntk_f()
		if err != nil {
			return err
		}
	}

	switch cv {
	case 82:
		{
			_, err = this.jj_consume_token(82)
			break
		}
	case 83:
		{
			_, err = this.jj_consume_token(83)
			break
		}
	case 84:
		{
			_, err = this.jj_consume_token(84)
			break
		}
	case 85:
		{
			_, err = this.jj_consume_token(85)
			break
		}
	case 86:
		{
			_, err = this.jj_consume_token(86)
			break
		}
	case 87:
		{
			_, err = this.jj_consume_token(87)
			break
		}
	case 88:
		{
			_, err = this.jj_consume_token(88)
			break
		}
	case 89:
		{
			_, err = this.jj_consume_token(89)
			break
		}
	case 90:
		{
			_, err = this.jj_consume_token(90)
			break
		}
	case 91:
		{
			_, err = this.jj_consume_token(91)
			break
		}
	case 92:
		{
			_, err = this.jj_consume_token(92)
			break
		}
	case 93:
		{
			_, err = this.jj_consume_token(93)
			break
		}
	case 94:
		{
			_, err = this.jj_consume_token(94)
			break
		}
	case AND:
		{
			_, err = this.jj_consume_token(AND)
			break
		}
	case OR:
		{
			_, err = this.jj_consume_token(OR)
			break
		}
	default:
		this.jj_consume_token(-1)
		err = newParseException("binop error")
	}
	return err
}

func (this *LuaParser) FieldSep() error { /*@bgen(jjtree) FieldSep */
	jjtn000 := NewSimpleNode(JJTFIELDSEP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()
	var err error
	cv := this.jj_ntk
	if cv == -1 {
		cv, err = this.jj_ntk_f()
		if err != nil {
			return err
		}
	}

	switch cv {
	case 72:
		{
			_, err = this.jj_consume_token(72)
			break
		}
	case 70:
		{
			_, err = this.jj_consume_token(70)
			break
		}
	default:
		this.jj_consume_token(-1)
		err = newParseException("FieldSep error")
	}
	return err
}

func (this *LuaParser) Field() error { /*@bgen(jjtree) Field */
	jjtn000 := NewSimpleNode(JJTFIELD)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 77:
			{
				_, err = this.jj_consume_token(77)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(78)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(71)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				break
			}
		default:
			var ok bool
			ok, err = this.jj_2_7(2)
			if err != nil {
				return err
			}
			if ok {
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(71)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
			} else {
				cv = this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}
				switch cv {
				case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, FALSE, FUNCTION,
					NIL, NOT, TRUE, NAME, NUMBER, STRING, CHARSTRING,
					69, 75, 79, 80, 83:
					{
						err = this.Exp()
						break
					}
				default:
					this.jj_consume_token(-1)
					err = newParseException("field error")
				}
			}
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FieldList() error { /*@bgen(jjtree) FieldList */
	jjtn000 := NewSimpleNode(JJTFIELDLIST)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		err = this.Field()
		if err != nil {
			return err
		}
		// label_9:
		for {
			var ok bool
			ok, err = this.jj_2_6(2)
			if err != nil {
				return err
			}
			if ok {

			} else {
				// break label_9
				break
			}
			err = this.FieldSep()
			if err != nil {
				return err
			}
			err = this.Field()
			if err != nil {
				return err
			}
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}
		switch cv {
		case 70, 72:
			{
				err = this.FieldSep()
				break
			}
		default:
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) TableConstructor() error { /*@bgen(jjtree) TableConstructor */
	jjtn000 := NewSimpleNode(JJTTABLECONSTRUCTOR)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(80)
		if err != nil {
			return err
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}
		switch cv {
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN,
			FALSE, FUNCTION, NIL, NOT, TRUE, NAME, NUMBER, STRING, CHARSTRING,
			69, 75, 77, 79, 80, 83:
			{
				err = this.FieldList()
				if err != nil {
					return err
				}
				break
			}
		default:
		}
		_, err = this.jj_consume_token(81)
		return err
	}()

	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) ParList() error { /*@bgen(jjtree) ParList */
	jjtn000 := NewSimpleNode(JJTPARLIST)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case NAME:
			{
				err = this.NameList()
				if err != nil {
					return err
				}
				cv = this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}
				switch cv {
				case 72:
					{
						_, err = this.jj_consume_token(72)
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(79)
						if err != nil {
							return err
						}
						break
					}
				default:
				}
				break
			}
		case 79:
			{
				_, err = this.jj_consume_token(79)
				break
			}
		default:
			this.jj_consume_token(-1)
			err = newParseException("ParList error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FuncBody() error { /*@bgen(jjtree) FuncBody */
	jjtn000 := NewSimpleNode(JJTFUNCBODY)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(75)
		if err != nil {
			return err
		}

		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case NAME, 79:
			{
				err = this.ParList()
				if err != nil {
					return err
				}
				break
			}
		default:
		}
		_, err = this.jj_consume_token(76)
		if err != nil {
			return err
		}
		err = this.Block()
		if err != nil {
			return err
		}
		_, err = this.jj_consume_token(END)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FunctionCall() error { /*@bgen(jjtree) FunctionCall */
	jjtn000 := NewSimpleNode(JJTFUNCTIONCALL)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(FUNCTION)
		if err != nil {
			return err
		}
		err = this.FuncBody()
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) SubExp() error { /*@bgen(jjtree) SubExp */
	jjtn000 := NewSimpleNode(JJTSUBEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN,
			FALSE, FUNCTION, NIL, TRUE, NAME, NUMBER, STRING, CHARSTRING,
			75, 79, 80:
			{
				err = this.SimpleExp()
				if err != nil {
					return err
				}
				break
			}
		case NOT, 69, 83:
			{
				err = this.Unop()
				if err != nil {
					return err
				}
				err = this.SubExp()
				if err != nil {
					return err
				}
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("subexp error")
		}
		// label_8:
		for {
			var ok bool
			ok, err = this.jj_2_5(2)
			if err != nil {
				return err
			}
			if ok {

			} else {
				// break label_8
				break
			}
			err = this.Binop()
			if err != nil {
				return err
			}
			err = this.SubExp()
			if err != nil {
				return err
			}
		}
		return nil
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Exp() error { /*@bgen(jjtree) Exp */
	jjtn000 := NewSimpleNode(JJTEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := this.SubExp()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Str() error { /*@bgen(jjtree) Str */
	jjtn000 := NewSimpleNode(JJTSTR)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case STRING:
			{
				_, err = this.jj_consume_token(STRING)
				break
			}
		case CHARSTRING:
			{
				_, err = this.jj_consume_token(CHARSTRING)
				break
			}
		case LONGSTRING0:
			{
				_, err = this.jj_consume_token(LONGSTRING0)
				break
			}
		case LONGSTRING1:
			{
				_, err = this.jj_consume_token(LONGSTRING1)
				break
			}
		case LONGSTRING2:
			{
				_, err = this.jj_consume_token(LONGSTRING2)
				break
			}
		case LONGSTRING3:
			{
				_, err = this.jj_consume_token(LONGSTRING3)
				break
			}
		case LONGSTRINGN:
			{
				_, err = this.jj_consume_token(LONGSTRINGN)
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("str error")
		}
		return err
	}()
	return err
}

func (this *LuaParser) SimpleExp() error { /*@bgen(jjtree) SimpleExp */
	jjtn000 := NewSimpleNode(JJTSIMPLEEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case NIL:
			{
				_, err = this.jj_consume_token(NIL)
				break
			}
		case TRUE:
			{
				_, err = this.jj_consume_token(TRUE)
				break
			}
		case FALSE:
			{
				_, err = this.jj_consume_token(FALSE)
				break
			}
		case NUMBER:
			{
				_, err = this.jj_consume_token(NUMBER)
				break
			}
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, STRING, CHARSTRING:
			{
				err = this.Str()
				break
			}
		case 79:
			{
				_, err = this.jj_consume_token(79)
				break
			}
		case 80:
			{
				err = this.TableConstructor()
				break
			}
		case FUNCTION:
			{
				err = this.FunctionCall()
				break
			}
		case NAME, 75:
			{
				_, err = this.PrimaryExp()
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("simpleExp error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) ExpList() error { /*@bgen(jjtree) ExpList */
	jjtn000 := NewSimpleNode(JJTEXPLIST)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		err = this.Exp()
		if err != nil {
			return err
		}
	label_7:
		for {
			cv := this.jj_ntk
			if cv == -1 {
				cv, err = this.jj_ntk_f()
				if err != nil {
					return err
				}
			}

			switch cv {
			case 72:
				{
					break
				}
			default:
				break label_7
			}
			_, err = this.jj_consume_token(72)
			if err != nil {
				return err
			}
			err = this.Exp()
			if err != nil {
				return err
			}
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) NameList() error { /*@bgen(jjtree) NameList */
	jjtn000 := NewSimpleNode(JJTNAMELIST)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(NAME)
		if err != nil {
			return err
		}
		// label_6:
		for {
			ok := false
			ok, err = this.jj_2_4(2)
			if err != nil {
				return err
			}
			if ok {

			} else {
				// break label_6
				break
			}
			_, err = this.jj_consume_token(72)
			if err != nil {
				return err
			}
			_, err = this.jj_consume_token(NAME)
			if err != nil {
				return err
			}
		}
		return err
	}()
	return err
}

func (this *LuaParser) FuncArgs() error { /*@bgen(jjtree) FuncArgs */
	jjtn000 := NewSimpleNode(JJTFUNCARGS)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 75:
			{
				_, err = this.jj_consume_token(75)
				if err != nil {
					return err
				}
				cv = this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}
				switch cv {
				case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN,
					FALSE, FUNCTION, NIL, NOT, TRUE, NAME, NUMBER, STRING, CHARSTRING,
					69, 75, 79, 80, 83:
					{
						err = this.ExpList()
						break
					}
				default:

				}
				_, err = this.jj_consume_token(76)
				if err != nil {
					return err
				}
				break
			}
		case 80:
			{
				err = this.TableConstructor()
				if err != nil {
					return err
				}
				break
			}
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, STRING, CHARSTRING:
			{
				err = this.Str()
				if err != nil {
					return err
				}
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("FuncArg error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FuncOp() error { /*@bgen(jjtree) FuncOp */
	jjtn000 := NewSimpleNode(JJTFUNCOP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		jjtn000.line = this.token.BeginLine

		switch cv {
		case 74:
			{
				_, err = this.jj_consume_token(74)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				err = this.FuncArgs()
				if err != nil {
					return err
				}
				break
			}
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, STRING, CHARSTRING,
			75, 80:
			{
				err = this.FuncArgs()
				if err != nil {
					return err
				}
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("FuncOp error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FieldOp() error { /*@bgen(jjtree) FieldOp */
	jjtn000 := NewSimpleNode(JJTFIELDOP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 73:
			{
				_, err = this.jj_consume_token(73)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				break
			}
		case 77:
			{
				_, err = this.jj_consume_token(77)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(78)
				if err != nil {
					return err
				}
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("FieldOp error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) PostfixOp() (int, error) { /*@bgen(jjtree) PostfixOp */
	jjtn000 := NewSimpleNode(JJTPOSTFIXOP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()
	r, err := func() (int, error) {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return 0, err
			}
		}

		switch cv {
		case 73, 77:
			{
				err = this.FieldOp()
				if err != nil {
					return 0, err
				}
				this.jjtree.closeNodeScopeB(jjtn000, true)
				jjtc000 = false
				// if "" != null {
				return VAR, nil
				// }
				// break
			}
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN, STRING, CHARSTRING,
			74, 75, 80:
			{
				err = this.FuncOp()
				if err != nil {
					return 0, nil
				}
				this.jjtree.closeNodeScopeB(jjtn000, true)
				jjtc000 = false
				//	if "" != null {
				return CALL, nil
				//	}
				//break
			}
		default:
			this.jj_consume_token(-1)
			return 0, newParseException("postfixOp error")
		}
		return -1, err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	if r == -1 {
		return 0, newParseException("Missing return statement in function")
	}
	return r, nil
}

func (this *LuaParser) PrimaryExp() (int, error) { /*@bgen(jjtree) PrimaryExp */
	jjtn000 := NewSimpleNode(JJTPRIMARYEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	r, err := func() (int, error) {
		var err error
		r := VAR
		err = this.PrefixExp()
		if err != nil {
			return 0, err
		}
		// label_5:
		for {
			ok := false
			ok, err = this.jj_2_3(2)
			if ok {
			} else {
				// break label_5;
				break
			}
			r, err = this.PostfixOp()
			if err != nil {
				return 0, err
			}
		}
		this.jjtree.closeNodeScopeB(jjtn000, true)
		jjtc000 = false
		// if "" != null {
		return r, nil
		// }
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	if r == -1 {
		return 0, newParseException("Missing return statement in function")
	}
	return r, err
}

func (this *LuaParser) ParenExp() error { /*@bgen(jjtree) ParenExp */
	jjtn000 := NewSimpleNode(JJTPARENEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(75)
		if err != nil {
			return err
		}
		err = this.Exp()
		if err != nil {
			return err
		}
		_, err = this.jj_consume_token(76)
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) PrefixExp() error { /*@bgen(jjtree) PrefixExp */
	jjtn000 := NewSimpleNode(JJTPREFIXEXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case NAME:
			{
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				break
			}
		case 75:
			{
				err = this.ParenExp()
				if err != nil {
					return err
				}
				break
			}
		default:
			this.jj_consume_token(-1)
			return newParseException("PrefixExp error")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) FuncName() error { /*@bgen(jjtree) FuncName */
	jjtn000 := NewSimpleNode(JJTFUNCNAME)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(NAME)
		if err != nil {
			return err
		}
	label_4:
		for {
			cv := this.jj_ntk
			if cv == -1 {
				cv, err = this.jj_ntk_f()
				if err != nil {
					return err
				}
			}

			switch cv {
			case 73:
				{
					break
				}
			default:
				break label_4
			}
			_, err = this.jj_consume_token(73)
			if err != nil {
				return err
			}
			_, err = this.jj_consume_token(NAME)
			return err
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 74:
			{
				_, err = this.jj_consume_token(74)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				break
			}
		default:

		}
		return err
	}()
	return err
}

func (this *LuaParser) VarExp() error { /*@bgen(jjtree) VarExp */
	jjtn000 := NewSimpleNode(JJTVAREXP)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		var r int
		r, err = this.PrimaryExp()
		if err != nil {
			return err
		}
		this.jjtree.closeNodeScopeB(jjtn000, true)
		jjtc000 = false
		if r != VAR {
			return newParseException("expected variable expression")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Assign() error { /*@bgen(jjtree) Assign */
	jjtn000 := NewSimpleNode(JJTASSIGN)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
	label_3:
		for {
			cv := this.jj_ntk
			if cv == -1 {
				cv, err = this.jj_ntk_f()
				if err != nil {
					return err
				}
			}

			switch cv {
			case 72:
				{
					break
				}
			default:
				break label_3
			}
			_, err = this.jj_consume_token(72)
			if err != nil {
				return err
			}
			err = this.VarExp()
			if err != nil {
				return err
			}
		}
		_, err = this.jj_consume_token(71)
		if err != nil {
			return err
		}
		err = this.ExpList()
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) ExprStat() error { /*@bgen(jjtree) ExprStat */
	jjtn000 := NewSimpleNode(JJTEXPRSTAT)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		var r int
		need := CALL

		r, err = this.PrimaryExp()
		if err != nil {
			return err
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 71, 72:
			{
				err = this.Assign()
				if err != nil {
					return err
				}
				need = VAR
				break
			}
		default:
		}
		this.jjtree.closeNodeScopeB(jjtn000, true)
		jjtc000 = false
		if r != need {
			return newParseException("expected function call or assignment")
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Label() error { /*@bgen(jjtree) Label */
	jjtn000 := NewSimpleNode(JJTLABEL)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(DBCOLON)
		if err != nil {
			return err
		}
		_, err = this.jj_consume_token(NAME)
		if err != nil {
			return err
		}
		_, err = this.jj_consume_token(DBCOLON)
		return err
	}()
	return err
}

func (this *LuaParser) ReturnStat() error { /*@bgen(jjtree) ReturnStat */
	jjtn000 := NewSimpleNode(JJTRETURNSTAT)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		_, err = this.jj_consume_token(RETURN)
		if err != nil {
			return err
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case LONGSTRING0, LONGSTRING1, LONGSTRING2, LONGSTRING3, LONGSTRINGN,
			FALSE, FUNCTION, NIL, NOT, TRUE, NAME, NUMBER, STRING, CHARSTRING,
			69, 75, 79, 80, 83:
			{
				err = this.ExpList()
				if err != nil {
					return err
				}
				break
			}
		default:
		}

		cv = this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 70:
			{
				_, err = this.jj_consume_token(70)
				if err != nil {
					return err
				}
				break
			}
		default:

		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Stat() error { /*@bgen(jjtree) Stat */
	jjtn000 := NewSimpleNode(JJTSTAT)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case 70:
			{
				_, err = this.jj_consume_token(70)
				if err != nil {
					return err
				}
				break
			}
		case DBCOLON:
			{
				err = this.Label()
				if err != nil {
					return err
				}
				break
			}
		case BREAK:
			{
				_, err = this.jj_consume_token(BREAK)
				if err != nil {
					return err
				}
				break
			}
		case GOTO:
			{
				_, err = this.jj_consume_token(GOTO)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				break
			}
		case DO:
			{
				_, err = this.jj_consume_token(DO)
				if err != nil {
					return err
				}
				err = this.Block()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(END)
				if err != nil {
					return err
				}
				break
			}
		case WHILE:
			{
				_, err = this.jj_consume_token(WHILE)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(DO)
				if err != nil {
					return err
				}
				err = this.Block()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(END)
				if err != nil {
					return err
				}
				break
			}
		case REPEAT:
			{
				_, err = this.jj_consume_token(REPEAT)
				if err != nil {
					return err
				}
				err = this.Block()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(UNTIL)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				break
			}
		case IF:
			{
				_, err = this.jj_consume_token(IF)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(THEN)
				if err != nil {
					return err
				}
				err = this.Block()
				if err != nil {
					return err
				}
			label_2:
				for {
					cv := this.jj_ntk
					if cv == -1 {
						cv, err = this.jj_ntk_f()
						if err != nil {
							return err
						}
					}

					switch cv {
					case ELSEIF:
						{
							break
						}
					default:
						break label_2
					}
					_, err = this.jj_consume_token(ELSEIF)
					if err != nil {
						return err
					}
					err = this.Exp()
					if err != nil {
						return err
					}
					_, err = this.jj_consume_token(THEN)
					if err != nil {
						return err
					}
					err = this.Block()
					if err != nil {
						return err
					}
				}
				cv := this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}

				switch cv {
				case ELSE:
					{
						_, err = this.jj_consume_token(ELSE)
						if err != nil {
							return err
						}
						err = this.Block()
						if err != nil {
							return err
						}
						break
					}
				default:

				}
				_, err = this.jj_consume_token(END)
				if err != nil {
					return err
				}
				break
			}
		default:
			ok := false
			ok, err = this.jj_2_1(3)
			if err != nil {
				return err
			}
			if ok {
				_, err = this.jj_consume_token(FOR)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(NAME)
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(71)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(72)
				if err != nil {
					return err
				}
				err = this.Exp()
				if err != nil {
					return err
				}

				cv := this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}

				switch cv {
				case 72:
					{
						_, err = this.jj_consume_token(72)
						if err != nil {
							return err
						}
						err = this.Exp()
						if err != nil {
							return err
						}
						break
					}
				default:

				}
				_, err = this.jj_consume_token(DO)
				if err != nil {
					return err
				}
				err = this.Block()
				if err != nil {
					return err
				}
				_, err = this.jj_consume_token(END)
				if err != nil {
					return err
				}
			} else {
				cv := this.jj_ntk
				if cv == -1 {
					cv, err = this.jj_ntk_f()
					if err != nil {
						return err
					}
				}

				switch cv {
				case FOR:
					{
						_, err = this.jj_consume_token(FOR)
						if err != nil {
							return err
						}
						err = this.NameList()
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(IN)
						if err != nil {
							return err
						}
						err = this.ExpList()
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(DO)
						if err != nil {
							return err
						}
						err = this.Block()
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(END)
						if err != nil {
							return err
						}
						break
					}
				case FUNCTION:
					{
						_, err = this.jj_consume_token(FUNCTION)
						if err != nil {
							return err
						}
						err = this.FuncName()
						if err != nil {
							return err
						}
						err = this.FuncBody()
						if err != nil {
							return err
						}
						break
					}
				default:
					ok := false
					ok, err = this.jj_2_2(2)
					if err != nil {
						return err
					}
					if ok {
						_, err = this.jj_consume_token(LOCAL)
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(FUNCTION)
						if err != nil {
							return err
						}
						_, err = this.jj_consume_token(NAME)
						if err != nil {
							return err
						}
						err = this.FuncBody()
						if err != nil {
							return err
						}
					} else {
						cv := this.jj_ntk
						if cv == -1 {
							cv, err = this.jj_ntk_f()
							if err != nil {
								return err
							}
						}

						switch cv {
						case LOCAL:
							{
								_, err = this.jj_consume_token(LOCAL)
								if err != nil {
									return err
								}
								err = this.NameList()
								if err != nil {
									return err
								}
								cv := this.jj_ntk
								if cv == -1 {
									cv, err = this.jj_ntk_f()
									if err != nil {
										return err
									}
								}

								switch cv {
								case 71:
									{
										_, err = this.jj_consume_token(71)
										if err != nil {
											return err
										}
										err = this.ExpList()
										if err != nil {
											return err
										}
										break
									}
								default:

								}
								break
							}
						case NAME, 75:
							{
								err = this.ExprStat()
								if err != nil {
									return err
								}
								break
							}
						default:
							_, err = this.jj_consume_token(-1)
							return newParseException("stat error")
						}
					}
				}
			}
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Block() error { /*@bgen(jjtree) Block */
	jjtn000 := NewSimpleNode(JJTBLOCK)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		// i := 0
	label_1:
		for {
			cv := this.jj_ntk
			if cv == -1 {
				cv, err = this.jj_ntk_f()
				if err != nil {
					return err
				}
			}

			switch cv {
			case BREAK, DO, FOR, FUNCTION, GOTO, IF, LOCAL, REPEAT,
				WHILE, NAME, DBCOLON, 70, 75:
				{
					break
				}
			default:
				break label_1
			}
			// i++
			// fmt.Println("block", i, cv)
			err = this.Stat()
			if err != nil {
				return err
			}
		}
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}

		switch cv {
		case RETURN:
			{
				err = this.ReturnStat()
				if err != nil {
					return err
				}
				break
			}
		default:

		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return err
}

func (this *LuaParser) Chunk() (Node, error) { /*@bgen(jjtree) Chunk */
	jjtn000 := NewSimpleNode(JJTCHUNK)
	jjtc000 := true
	this.jjtree.openNodeScope(jjtn000)
	defer func() {
		if jjtc000 {
			this.jjtree.closeNodeScopeB(jjtn000, true)
		}
	}()

	err := func() error {
		var err error
		cv := this.jj_ntk
		if cv == -1 {
			cv, err = this.jj_ntk_f()
			if err != nil {
				return err
			}
		}
		switch cv {
		case 69:
			{
				_, err = this.jj_consume_token(69)
				if err != nil {
					return err
				}
				this.token_source.SwitchTo(IN_COMMENT)
				break
			}
		default:

		}
		err = this.Block()
		if err != nil {
			return err
		}
		_, err = this.jj_consume_token(0)
		if err != nil {
			return err
		}
		return err
	}()
	if err != nil {
		if jjtc000 {
			this.jjtree.clearNodeScope(jjtn000)
			jjtc000 = false
		} else {
			this.jjtree.popNode()
		}
	}
	return jjtn000, err
}
