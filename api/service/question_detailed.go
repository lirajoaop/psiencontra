package service

var detailedQuestions = []Question{
	// ===== BLOCK 1: THEORETICAL APPROACHES — LIKERT (40 items: 4 per approach + 1 reverse per approach) =====
	// Each statement is rated 1-5 (Discordo totalmente → Concordo totalmente).
	// Items are shuffled at runtime within each block to prevent halo/carry-over
	// effects and reduce construct transparency.
	// Dimensions: A=visão de sujeito, B=sofrimento/mudança, C=postura clínica, D=epistemologia

	// --- Psicanálise ---
	{
		ID:      101,
		Text:    "Grande parte do que sentimos e fazemos é determinado por forças inconscientes que não controlamos diretamente.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("psicanalise"),
	},
	{
		ID:      102,
		Text:    "O sofrimento psíquico frequentemente tem raízes em experiências infantis e conflitos que ficaram sem elaboração.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("psicanalise"),
	},
	{
		ID:      103,
		Text:    "A escuta do que o paciente diz — e especialmente do que ele evita dizer — é a ferramenta mais importante da clínica.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("psicanalise"),
	},
	{
		ID:      104,
		Text:    "A análise de sonhos, lapsos e associações livres revela verdades profundas sobre a vida psíquica.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("psicanalise"),
	},

	// --- Fenomenologia-Existencial ---
	{
		ID:      105,
		Text:    "Cada pessoa é livre e responsável por dar sentido à própria existência, mesmo diante do sofrimento.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("fenomenologia"),
	},
	{
		ID:      106,
		Text:    "A angústia diante da finitude e das escolhas não é patologia — é parte essencial da condição humana.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("fenomenologia"),
	},
	{
		ID:      107,
		Text:    "O terapeuta deve suspender suas teorias e preconceitos para acessar o mundo tal como o paciente o vivencia.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("fenomenologia"),
	},
	{
		ID:      108,
		Text:    "Compreender como a pessoa vivencia sua experiência é mais importante do que classificar ou explicar seu comportamento.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("fenomenologia"),
	},

	// --- Análise do Comportamento ---
	{
		ID:      109,
		Text:    "O comportamento humano é moldado principalmente pelas consequências que recebe do ambiente ao longo da vida.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("comportamental"),
	},
	{
		ID:      110,
		Text:    "Mudança terapêutica efetiva acontece quando identificamos e modificamos as contingências que mantêm o problema.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("comportamental"),
	},
	{
		ID:      111,
		Text:    "Intervenções com metas claras, mensuráveis e baseadas em observação sistemática são mais eficazes que abordagens intuitivas.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("comportamental"),
	},
	{
		ID:      112,
		Text:    "A psicologia deve se pautar por dados empíricos e observação direta, não por conceitos que não podem ser verificados.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("comportamental"),
	},

	// --- Terapia Cognitivo-Comportamental ---
	{
		ID:      113,
		Text:    "A forma como interpretamos os eventos tem mais impacto nas nossas emoções do que os eventos em si.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("tcc"),
	},
	{
		ID:      114,
		Text:    "Identificar e reestruturar pensamentos disfuncionais é uma das formas mais eficazes de reduzir o sofrimento.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("tcc"),
	},
	{
		ID:      115,
		Text:    "Para mim, terapeuta e paciente devem trabalhar colaborativamente: identificar demandas concretas, aplicar técnicas direcionadas e verificar juntos se o tratamento está gerando mudança real.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("tcc"),
	},
	{
		ID:      116,
		Text:    "Tratamentos com objetivos claros por sessão e uso de escalas para acompanhar a evolução do paciente me parecem uma boa forma de fazer terapia.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("tcc"),
	},

	// --- Psicologia Analítica (Jung) ---
	{
		ID:      117,
		Text:    "O ser humano busca naturalmente integrar aspectos opostos de si mesmo — consciência e inconsciente, persona e sombra.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("junguiana"),
	},
	{
		ID:      118,
		Text:    "O objetivo da terapia não é apenas aliviar sintomas, mas ajudar a pessoa a se tornar mais completa: integrando luz e sombra, razão e emoção.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("junguiana"),
	},
	{
		ID:      119,
		Text:    "Sonhos, mitos, contos e símbolos são linguagens privilegiadas do inconsciente e devem ser explorados na clínica.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("junguiana"),
	},
	{
		ID:      120,
		Text:    "Existe uma dimensão coletiva do inconsciente, com arquétipos e padrões compartilhados por toda a humanidade.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("junguiana"),
	},

	// --- Gestalt-terapia ---
	{
		ID:      121,
		Text:    "O ser humano é uma totalidade — corpo, emoções, pensamentos e relações formam um todo inseparável.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("gestalt"),
	},
	{
		ID:      122,
		Text:    "Muitas formas de sofrimento nascem da interrupção do contato genuíno consigo mesmo e com o ambiente.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("gestalt"),
	},
	{
		ID:      123,
		Text:    "O que está acontecendo aqui e agora na sessão tem mais valor terapêutico do que reconstruir o passado.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("gestalt"),
	},
	{
		ID:      124,
		Text:    "A experiência vivida e percebida no momento presente é mais reveladora do que qualquer explicação teórica.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("gestalt"),
	},

	// --- Psicologia Sócio-Histórica ---
	{
		ID:      125,
		Text:    "O ser humano se constitui nas relações sociais — somos, ao mesmo tempo, produto e produtores da nossa cultura.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("socio_historica"),
	},
	{
		ID:      126,
		Text:    "O sofrimento psíquico não pode ser entendido sem considerar as condições sociais, econômicas e políticas da pessoa.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("socio_historica"),
	},
	{
		ID:      127,
		Text:    "O psicólogo tem responsabilidade social: deve atuar na defesa de direitos e na construção de políticas públicas.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("socio_historica"),
	},
	{
		ID:      128,
		Text:    "A psicologia comprometida com a transformação das desigualdades sociais é mais relevante do que a focada apenas no indivíduo.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("socio_historica"),
	},

	// --- Humanismo / Abordagem Centrada na Pessoa ---
	{
		ID:      129,
		Text:    "Todo ser humano possui uma tendência inata ao crescimento, à autonomia e à realização do seu potencial.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("humanismo"),
	},
	{
		ID:      130,
		Text:    "O sofrimento surge quando a pessoa precisa negar partes de sua experiência para manter a aceitação dos outros.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("humanismo"),
	},
	{
		ID:      131,
		Text:    "Empatia genuína, aceitação incondicional e autenticidade do terapeuta são os ingredientes mais importantes da mudança terapêutica.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("humanismo"),
	},
	{
		ID:      132,
		Text:    "Confio na capacidade de cada pessoa de encontrar suas próprias respostas quando é escutada sem julgamento, mesmo em crises profundas.",
		Type:    "likert",
		Block:   "approaches",
		Options: likertOptions("humanismo"),
	},

	// ===== REVERSE ITEMS (1 per approach, 8 total) =====
	// These are phrased *against* the construct. The scoring service inverts
	// the rating (6 - r) before ipsativization. Reverse items reduce item
	// transparency and flag acquiescent/inattentive responders whose profile
	// becomes internally incoherent (high on tese AND anti-tese simultaneously).

	{
		ID:       133,
		Text:     "A maior parte do sofrimento psíquico pode ser explicada sem recorrer a conceitos como inconsciente ou conflitos infantis.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("psicanalise"),
		Reversed: true,
	},
	{
		ID:       134,
		Text:     "Explicar o comportamento humano por leis gerais e categorias objetivas é mais útil do que se deter na experiência singular de cada pessoa.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("fenomenologia"),
		Reversed: true,
	},
	{
		ID:       135,
		Text:     "Conceitos sobre processos internos — desejos, significados, sentidos — são indispensáveis na clínica, mesmo quando não podem ser observados diretamente.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("comportamental"),
		Reversed: true,
	},
	{
		ID:       136,
		Text:     "Protocolos estruturados e técnicas padronizadas tendem a engessar a terapia; o trabalho clínico deveria ser mais aberto e menos manualizado.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("tcc"),
		Reversed: true,
	},
	{
		ID:       137,
		Text:     "Arquétipos e inconsciente coletivo são metáforas literárias interessantes, mas têm pouco valor explicativo para a prática clínica.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("junguiana"),
		Reversed: true,
	},
	{
		ID:       138,
		Text:     "Para entender o sofrimento, é mais produtivo reconstruir a história do paciente do que insistir no que se passa no aqui-e-agora.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("gestalt"),
		Reversed: true,
	},
	{
		ID:       139,
		Text:     "O foco do psicólogo deve ser o indivíduo e sua vida interna; questões sociais e políticas competem a outros profissionais.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("socio_historica"),
		Reversed: true,
	},
	{
		ID:       140,
		Text:     "Sem técnicas específicas e bem definidas, a relação terapêutica sozinha dificilmente gera mudanças clínicas consistentes.",
		Type:     "likert",
		Block:    "approaches",
		Options:  likertOptions("humanismo"),
		Reversed: true,
	},

	// ===== BLOCK 2: FIELDS OF PRACTICE — LIKERT (36 items, 4 per field) =====
	// Dimensions: A=contexto/público, B=atividade/intervenção, C=motivação

	// --- Psicologia Clínica ---
	{
		ID:      201,
		Text:    "Me imagino atendendo pacientes em consultório, ajudando pessoas a atravessar crises emocionais e processos de autoconhecimento.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("clinica"),
	},
	{
		ID:      202,
		Text:    "Conduzir psicoterapia — individual, de casal ou em grupo — é a atividade que mais me atrai na psicologia.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("clinica"),
	},
	{
		ID:      203,
		Text:    "Minha principal motivação é acolher o sofrimento humano e ajudar as pessoas a encontrarem caminhos de melhora.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("clinica"),
	},
	{
		ID:      228,
		Text:    "Construir vínculo terapêutico de longo prazo, acompanhando o processo de transformação de um mesmo paciente ao longo do tempo, é algo com que me identifico.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("clinica"),
	},

	// --- Psicologia Organizacional ---
	{
		ID:      204,
		Text:    "Me interesso por compreender como fatores psicológicos afetam o desempenho, a satisfação e o bem-estar no trabalho.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("organizacional"),
	},
	{
		ID:      205,
		Text:    "Atividades como recrutamento, desenvolvimento de lideranças e gestão de clima organizacional me atraem.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("organizacional"),
	},
	{
		ID:      206,
		Text:    "Gostaria de usar a psicologia para tornar organizações mais saudáveis, produtivas e humanas.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("organizacional"),
	},
	{
		ID:      229,
		Text:    "Conduzir avaliações psicológicas em processos seletivos e orientar programas de desenvolvimento de carreira são atividades que me atraem.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("organizacional"),
	},

	// --- Psicologia Escolar/Educacional ---
	{
		ID:      207,
		Text:    "Atuar em escolas ou universidades, apoiando alunos, famílias e educadores, é algo que me atrai.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("escolar"),
	},
	{
		ID:      208,
		Text:    "Me interesso por orientação educacional, inclusão, dificuldades de aprendizagem e desenvolvimento infantojuvenil.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("escolar"),
	},
	{
		ID:      209,
		Text:    "Acredito que a psicologia pode transformar a educação e promover condições melhores de aprendizagem para todos.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("escolar"),
	},
	{
		ID:      230,
		Text:    "Me imagino mediando conflitos entre alunos, professores e famílias e estruturando ações preventivas de saúde mental na escola.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("escolar"),
	},

	// --- Psicologia Social e Comunitária ---
	{
		ID:      210,
		Text:    "Me sinto atraído(a) por trabalhar com comunidades em situação de vulnerabilidade, em CRAS, CREAS ou ONGs.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("social"),
	},
	{
		ID:      211,
		Text:    "Mobilização comunitária, facilitação de grupos e construção de políticas públicas de saúde mental me interessam.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("social"),
	},
	{
		ID:      212,
		Text:    "A psicologia deve ser instrumento de transformação social e defesa dos direitos de populações marginalizadas.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("social"),
	},
	{
		ID:      231,
		Text:    "Acompanhar famílias atendidas pelo SUAS, articulando rede com assistência social e saúde, é um trabalho com que me identifico.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("social"),
	},

	// --- Psicologia da Saúde/Hospitalar ---
	{
		ID:      213,
		Text:    "Me imagino atuando em hospitais, UBS ou equipes multidisciplinares, cuidando de pacientes em tratamento de saúde.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("saude"),
	},
	{
		ID:      214,
		Text:    "Acompanhamento psicológico hospitalar, preparação para cirurgias e cuidados paliativos são atividades que me atraem.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("saude"),
	},
	{
		ID:      215,
		Text:    "Quero contribuir para o cuidado integral, integrando saúde mental e saúde física no atendimento ao paciente.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("saude"),
	},
	{
		ID:      232,
		Text:    "Me atrai integrar equipes multidisciplinares em UBS, trabalhando promoção de saúde e atenção psicossocial junto ao SUS.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("saude"),
	},

	// --- Psicologia Jurídica ---
	{
		ID:      216,
		Text:    "A interface entre psicologia e sistema de justiça — varas de família, sistema penitenciário, medidas socioeducativas — me atrai.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("juridica"),
	},
	{
		ID:      217,
		Text:    "Elaboração de laudos periciais, avaliação de risco, mediação de conflitos e acompanhamento de vítimas me interessam.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("juridica"),
	},
	{
		ID:      218,
		Text:    "Contribuir para decisões judiciais mais humanizadas e para a garantia de direitos é algo que me motiva.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("juridica"),
	},
	{
		ID:      233,
		Text:    "Atuar em Varas da Infância, medidas de proteção e processos de adoção, escutando crianças e famílias em situações complexas, é uma perspectiva que me interessa.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("juridica"),
	},

	// --- Psicologia do Esporte ---
	{
		ID:      219,
		Text:    "Me atrai trabalhar com atletas e equipes esportivas, do amador ao alto rendimento.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("esporte"),
	},
	{
		ID:      220,
		Text:    "Preparação mental, controle de ansiedade competitiva, foco atencional e recuperação pós-lesão me interessam.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("esporte"),
	},
	{
		ID:      221,
		Text:    "Gostaria de ajudar atletas a alcançarem seu potencial máximo, trabalhando os aspectos psicológicos da performance.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("esporte"),
	},
	{
		ID:      234,
		Text:    "Acompanhar atletas em transições de carreira — lesões, aposentadoria, pressão midiática — é um aspecto da psicologia do esporte que me atrai.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("esporte"),
	},

	// --- Neuropsicologia ---
	{
		ID:      222,
		Text:    "Me fascina a relação entre cérebro e comportamento — como processos neurológicos influenciam cognição, emoção e personalidade.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("neuropsicologia"),
	},
	{
		ID:      223,
		Text:    "Avaliação neuropsicológica, aplicação de baterias de testes cognitivos e programas de reabilitação me atraem.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("neuropsicologia"),
	},
	{
		ID:      224,
		Text:    "Gostaria de atuar em centros de neurociências, clínicas de reabilitação ou hospitais, trabalhando com funções cognitivas.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("neuropsicologia"),
	},
	{
		ID:      235,
		Text:    "Me atrai trabalhar com pacientes pós-AVC, TCE, TDAH ou quadros de declínio cognitivo, avaliando e reabilitando funções cognitivas.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("neuropsicologia"),
	},

	// --- Psicometria ---
	{
		ID:      225,
		Text:    "Me interesso pela ciência por trás dos testes psicológicos — como são construídos, validados e interpretados.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("psicometria"),
	},
	{
		ID:      226,
		Text:    "Análise estatística de dados, construção de escalas e pesquisa em avaliação psicológica me atraem.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("psicometria"),
	},
	{
		ID:      227,
		Text:    "Acredito que a psicologia precisa de instrumentos de medida cada vez mais precisos, justos e cientificamente fundamentados.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("psicometria"),
	},
	{
		ID:      236,
		Text:    "Me interessaria contribuir para a validação e adaptação de instrumentos de avaliação psicológica para a realidade brasileira.",
		Type:    "likert",
		Block:   "fields",
		Options: likertOptions("psicometria"),
	},

}

// likertOptions returns the standard 1-5 Likert scale options, all mapping to
// the given approach or field key. The frontend renders these as a scale, not
// as individual buttons.
func likertOptions(mapping string) []Option {
	return []Option{
		{Label: "Discordo totalmente", Value: "1", Mapping: mapping},
		{Label: "Discordo", Value: "2", Mapping: mapping},
		{Label: "Neutro", Value: "3", Mapping: mapping},
		{Label: "Concordo", Value: "4", Mapping: mapping},
		{Label: "Concordo totalmente", Value: "5", Mapping: mapping},
	}
}
