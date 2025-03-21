package shared

type experienceTable struct {
	Level      int
	XPNeeded   uint64
	XPPerLevel uint64
}

type ExperienceTable []experienceTable

func BuildExperienceTable() *ExperienceTable {
	ExperienceTable := make(ExperienceTable, 0, 101)

	for i := 0; i <= 100; i++ {
		ExperienceTable = append(ExperienceTable, experienceTable{
			Level:      i,
			XPNeeded:   xpNeeded(i),
			XPPerLevel: xpPerLevel(i),
		})
	}

	return &ExperienceTable
}

func (E *ExperienceTable) Level(xp uint64) int {
	for i, t := range *E {
		if xp < t.XPNeeded {
			return i
		}
	}

	return len(*E)
}

/*
// Private Functions
*/

func xpNeeded(level int) uint64 {
	if level == 0 {
		return 2400
	}

	if level == 1 {
		return 6210
	}

	return xpNeeded(level-1) + xpPerLevel(level-1)
}

func xpPerLevel(level int) uint64 {
	switch true {
	case level == 1:
		return 6210
	case level == 2:
		return 10120
	case level == 3:
		return 13800
	case level == 4:
		return 16100
	case level == 5:
		return 19200
	case level <= 7:
		return 21600
	case level == 8:
		return 24000
	case level == 9:
		return 26470
	case level < 15:
		return 32900
	case level < 20:
		return 54800
	case level < 25:
		return 75000
	case level < 30:
		return 85000
	case level < 35:
		return 108400
	case level < 40:
		return 121600
	case level < 45:
		return 130000
	case level < 50:
		return 138000
	default:
		return 144000
	}
}

/*
0	0	2,400	0	0	0	0
1	2,400	6,210	1	1	4	4
2	8,610	10,120	2	3	3	7
3	18,730	13,800	1	4	4	11
4	32,530	16,100	2	6	3	14
5	48,630	19,200	1	7	4	18
6	67,830	21,600	2	9	3	21
7	89,430	21,600	1	10	4	25
8	111,030	24,000	2	12	3	28
9	135,030	26,470	1	13	4	32
10	161,500	32,900	2	15	3	35

11	194,400	32,900	1	16	4	39
12	227,300	32,900	2	18	3	42
13	260,200	32,900	1	19	4	46
14	293,100	32,900	2	21	3	49

15	326,000	54,800	1	22	4	53
16	380,800	54,800	2	24	3	56
17	435,600	54,800	1	25	4	60
18	490,400	54,800	2	27	3	63
19	545,200	54,800	1	28	4	67

20	600,000	75,000	2	30	3	70

25	975,000	85,000	7	37	18	88
30	1,400,000	108,400	8	45	17	105
35	1,942,000	121,600	7	52	18	123
40	2,550,000	130,000	8	60	17	140
45	3,200,000	138,000	7	67	18	158
50	3,890,000	144,000	8	75	17	175
55	4,610,000	144,000	7	82	18	193
60	5,330,000	144,000	8	90	17	210
65	6,050,000	144,000	0	90	18	228
70	6,770,000	144,000	0	90	17	245
75	7,490,000	144,000	0	90	18	263
80	8,210,000	144,000	0	90	17	280
85	8,930,000	144,000	0	90	18	298
90	9,650,000	144,000	0	90	17	315
95	10,370,000	144,000	0	90	18	333
100	11,090,000	144,000
*/
