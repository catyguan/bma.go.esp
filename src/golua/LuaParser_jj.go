package golua

func (this *LuaParser) jj_2_1(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_1()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_2(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_2()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_3(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_3()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_4(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_4()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_5(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_5()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_6(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_6()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_2_7(xla int) (bool, error) {
	this.jj_la = xla
	this.jj_scanpos = this.token
	this.jj_lastpos = this.token
	ok, err := this.jj_3_7()
	if err != nil {
		if _, la := err.(*lookaheadSuccess); la {
			return true, nil
		}
		return false, err
	}
	return !ok, nil
}

func (this *LuaParser) jj_3R_23() (bool, error) {
	var ok bool
	var err error
	xsp := this.jj_scanpos
	ok, err = this.jj_scan_token(42)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(48)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(35)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(52)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_30()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(79)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_31()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_32()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_33()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3_4() (bool, error) {
	var ok bool
	var err error
	ok, err = this.jj_scan_token(72)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	ok, err = this.jj_scan_token(NAME)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_45() (bool, error) {
	var ok bool
	var err error
	ok, err = this.jj_3R_25()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_24() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_scan_token(83)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(43)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(69)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil

}

func (this *LuaParser) jj_3_2() (bool, error) {
	var ok bool
	var err error
	ok, err = this.jj_scan_token(LOCAL)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(FUNCTION)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_43() (bool, error) {
	var ok bool
	var err error
	ok, err = this.jj_3R_45()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3_1() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(FOR)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(NAME)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(71)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_41() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_35()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_11() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_scan_token(82)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(83)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(84)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(85)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(86)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(87)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(88)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(89)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(90)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(91)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(92)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(93)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(94)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(29)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(44)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_40() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_36()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_13() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_scan_token(72)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(70)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return true, nil
}

func (this *LuaParser) jj_3R_39() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(75)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	xsp := this.jj_scanpos

	ok, err = this.jj_3R_43()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(76)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_34() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_39()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_40()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_41()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3_6() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_13()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_3R_14()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_20() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_25()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_29() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_34()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3_7() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(NAME)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(71)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_22() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_28()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_29()
	if err != nil {
		return false, nil
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_28() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(74)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(NAME)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_19() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(77)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_14() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_19()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3_7()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_20()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	return true, nil
}

func (this *LuaParser) jj_3_5() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_11()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_3R_12()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_48() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_14()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_46() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_48()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3_3() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_10()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_27() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(77)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_3R_25()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_21() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_26()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_27()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_26() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(73)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	ok, err = this.jj_scan_token(NAME)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	return false, nil
}

func (this *LuaParser) jj_3R_36() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(80)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}

	xsp := this.jj_scanpos

	ok, err = this.jj_3R_46()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(81)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_16() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_22()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_15() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_21()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_10() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_15()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_16()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_38() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_42()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_18() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_24()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_47() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(75)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_37() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_scan_token(FUNCTION)
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_44() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_47()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_42() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_scan_token(51)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_44()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return true, nil
}

func (this *LuaParser) jj_3R_17() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_23()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_12() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_3R_17()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_3R_18()
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_25() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_12()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_35() (bool, error) {
	xsp := this.jj_scanpos

	var ok bool
	var err error

	ok, err = this.jj_scan_token(61)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(62)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(23)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(24)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(25)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(26)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	this.jj_scanpos = xsp

	ok, err = this.jj_scan_token(27)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	return true, nil
}

func (this *LuaParser) jj_3R_33() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_38()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_32() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_37()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_31() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_36()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}

func (this *LuaParser) jj_3R_30() (bool, error) {
	var ok bool
	var err error

	ok, err = this.jj_3R_35()
	if err != nil {
		return false, err
	}
	if ok {
		return true, nil
	}
	return false, nil
}
