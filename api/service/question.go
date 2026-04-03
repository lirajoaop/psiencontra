package service

type Option struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Mapping string `json:"mapping,omitempty"`
}

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"` // multiple_choice | open_ended
	Block   string   `json:"block"` // approaches | fields
	Options []Option `json:"options,omitempty"`
}

type QuestionService struct{}

func NewQuestionService() *QuestionService {
	return &QuestionService{}
}

func (s *QuestionService) GetAll() []Question {
	return questions
}

var questions = []Question{
	// ===== BLOCO 1: ABORDAGENS TEORICAS (10 perguntas) =====
	{
		ID:    1,
		Text:  "Na sua opiniao, qual e a principal causa do sofrimento humano?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Conflitos inconscientes nao resolvidos", Value: "conflitos_inconscientes", Mapping: "psicanalise"},
			{Label: "A perda de sentido e a dificuldade de ser autentico", Value: "perda_sentido", Mapping: "fenomenologia"},
			{Label: "Padroes de comportamento aprendidos que geram consequencias negativas", Value: "padroes_comportamento", Mapping: "comportamental"},
			{Label: "Pensamentos distorcidos sobre si mesmo e o mundo", Value: "pensamentos_distorcidos", Mapping: "tcc"},
			{Label: "O desequilibrio entre aspectos conscientes e inconscientes da personalidade", Value: "desequilibrio_personalidade", Mapping: "junguiana"},
			{Label: "A interrupcao do contato consigo mesmo e com o ambiente", Value: "interrupcao_contato", Mapping: "gestalt"},
			{Label: "As condicoes sociais e historicas que oprimem o individuo", Value: "condicoes_sociais", Mapping: "socio_historica"},
			{Label: "Disfuncoes nos padroes de comunicacao e relacionamento", Value: "disfuncoes_comunicacao", Mapping: "sistemica"},
		},
	},
	{
		ID:    2,
		Text:  "Qual visao sobre a natureza humana mais ressoa com voce?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Somos movidos por desejos inconscientes e pulsoes", Value: "desejos_inconscientes", Mapping: "psicanalise"},
			{Label: "Somos seres livres, responsaveis por criar nosso proprio sentido", Value: "seres_livres", Mapping: "fenomenologia"},
			{Label: "Somos produto do nosso ambiente e das consequencias dos nossos atos", Value: "produto_ambiente", Mapping: "comportamental"},
			{Label: "Temos capacidade racional de reestruturar nossos pensamentos", Value: "capacidade_racional", Mapping: "tcc"},
			{Label: "Buscamos a individuacao — integrar luz e sombra", Value: "individuacao", Mapping: "junguiana"},
			{Label: "Somos totalidades em constante processo de ajustamento criativo", Value: "totalidades", Mapping: "gestalt"},
			{Label: "Somos seres historicos, construidos nas relacoes sociais", Value: "seres_historicos", Mapping: "socio_historica"},
			{Label: "Somos parte de sistemas interconectados que nos influenciam", Value: "sistemas_interconectados", Mapping: "sistemica"},
		},
	},
	{
		ID:    3,
		Text:  "Em poucas palavras, o que significa 'cura' ou 'melhora' em psicologia para voce?",
		Type:  "open_ended",
		Block: "approaches",
	},
	{
		ID:    4,
		Text:  "Qual deveria ser o papel principal do terapeuta?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Interpretar o inconsciente e os significados ocultos", Value: "interpretar_inconsciente", Mapping: "psicanalise"},
			{Label: "Estar presente e facilitar a exploracao da experiencia vivida", Value: "estar_presente", Mapping: "fenomenologia"},
			{Label: "Analisar e modificar comportamentos atraves de tecnicas baseadas em evidencias", Value: "analisar_comportamentos", Mapping: "comportamental"},
			{Label: "Ajudar a identificar e reestruturar pensamentos disfuncionais", Value: "reestruturar_pensamentos", Mapping: "tcc"},
			{Label: "Guiar o paciente na exploracao de simbolos e arquetipos", Value: "explorar_simbolos", Mapping: "junguiana"},
			{Label: "Ampliar a awareness e o contato com o momento presente", Value: "ampliar_awareness", Mapping: "gestalt"},
			{Label: "Promover consciencia critica e transformacao social", Value: "consciencia_critica", Mapping: "socio_historica"},
			{Label: "Entender e reorganizar os padroes relacionais do sistema", Value: "reorganizar_padroes", Mapping: "sistemica"},
		},
	},
	{
		ID:    5,
		Text:  "Qual dessas afirmacoes mais ressoa com voce?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "\"O sonho e a estrada real para o inconsciente\"", Value: "sonho_inconsciente", Mapping: "psicanalise"},
			{Label: "\"Tornar-se pessoa e o objetivo da vida\"", Value: "tornar_pessoa", Mapping: "fenomenologia"},
			{Label: "\"O comportamento e funcao de suas consequencias\"", Value: "comportamento_consequencias", Mapping: "comportamental"},
			{Label: "\"Nao sao as coisas que nos perturbam, mas a visao que temos delas\"", Value: "visao_coisas", Mapping: "tcc"},
			{Label: "\"Quem olha para fora sonha; quem olha para dentro desperta\"", Value: "olha_dentro", Mapping: "junguiana"},
			{Label: "\"Perder a cabeca e cair em si\"", Value: "cair_em_si", Mapping: "gestalt"},
			{Label: "\"O homem se faz ao fazer o mundo e o mundo o faz\"", Value: "homem_mundo", Mapping: "socio_historica"},
			{Label: "\"O problema nao esta na pessoa, esta no padrao\"", Value: "problema_padrao", Mapping: "sistemica"},
		},
	},
	{
		ID:    6,
		Text:  "Se voce pudesse usar uma unica palavra para definir o objetivo da psicoterapia, qual seria?",
		Type:  "open_ended",
		Block: "approaches",
	},
	{
		ID:    7,
		Text:  "Que tipo de evidencia voce mais valoriza na psicologia?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "A escuta profunda do discurso e dos lapsos do paciente", Value: "escuta_profunda", Mapping: "psicanalise"},
			{Label: "A descricao detalhada da experiencia vivida", Value: "experiencia_vivida", Mapping: "fenomenologia"},
			{Label: "Dados experimentais e observacoes sistematicas do comportamento", Value: "dados_experimentais", Mapping: "comportamental"},
			{Label: "Estudos clinicos controlados e protocolos validados", Value: "estudos_controlados", Mapping: "tcc"},
			{Label: "A analise de sonhos, mitos e simbolos universais", Value: "sonhos_mitos", Mapping: "junguiana"},
			{Label: "A observacao fenomenologica do aqui-e-agora", Value: "aqui_agora", Mapping: "gestalt"},
			{Label: "A analise do contexto social, politico e historico", Value: "contexto_social", Mapping: "socio_historica"},
			{Label: "A observacao dos padroes de interacao no sistema familiar", Value: "padroes_interacao", Mapping: "sistemica"},
		},
	},
	{
		ID:    8,
		Text:  "Qual dessas atividades profissionais mais te atrairia?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Conduzir sessoes de analise com associacao livre e interpretacao de sonhos", Value: "analise_sonhos", Mapping: "psicanalise"},
			{Label: "Acompanhar alguem em uma crise existencial com escuta empática", Value: "crise_existencial", Mapping: "fenomenologia"},
			{Label: "Desenvolver um programa de modificacao de comportamento", Value: "modificacao_comportamento", Mapping: "comportamental"},
			{Label: "Aplicar tecnicas estruturadas para tratar ansiedade ou depressao", Value: "tecnicas_estruturadas", Mapping: "tcc"},
			{Label: "Facilitar um processo de autoconhecimento profundo usando simbolos", Value: "autoconhecimento_simbolos", Mapping: "junguiana"},
			{Label: "Conduzir um grupo de vivencia focado na percepcao corporal", Value: "vivencia_corporal", Mapping: "gestalt"},
			{Label: "Coordenar um projeto de intervencao comunitaria", Value: "intervencao_comunitaria", Mapping: "socio_historica"},
			{Label: "Fazer terapia familiar trabalhando as dinamicas do sistema", Value: "terapia_familiar", Mapping: "sistemica"},
		},
	},
	{
		ID:    9,
		Text:  "Qual filosofo, pensador ou autor mais influenciou (ou influenciaria) sua visao de mundo?",
		Type:  "open_ended",
		Block: "approaches",
	},
	{
		ID:    10,
		Text:  "Um paciente diz: 'Nao sei mais quem eu sou.' Como voce abordaria essa fala?",
		Type:  "open_ended",
		Block: "approaches",
	},
	// ===== BLOCO 2: CAMPOS DE ATUACAO (5 perguntas) =====
	{
		ID:    11,
		Text:  "Com qual publico voce mais gostaria de trabalhar?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Criancas e adolescentes em sofrimento emocional", Value: "criancas_adolescentes", Mapping: "clinica"},
			{Label: "Adultos em crises pessoais e emocionais", Value: "adultos_crises", Mapping: "clinica"},
			{Label: "Atletas e esportistas", Value: "atletas", Mapping: "esporte"},
			{Label: "Funcionarios e gestores de empresas", Value: "funcionarios_empresas", Mapping: "organizacional"},
			{Label: "Comunidades em situacao de vulnerabilidade", Value: "comunidades", Mapping: "social"},
			{Label: "Pacientes hospitalizados ou em tratamento medico", Value: "pacientes_hospitalizados", Mapping: "saude"},
			{Label: "Pessoas envolvidas em processos judiciais ou em situacao carceraria", Value: "processos_judiciais", Mapping: "juridica"},
			{Label: "Alunos e educadores em contexto escolar", Value: "alunos_educadores", Mapping: "escolar"},
		},
	},
	{
		ID:    12,
		Text:  "Qual ambiente de trabalho mais combina com voce?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Consultorio particular", Value: "consultorio", Mapping: "clinica"},
			{Label: "Empresa ou departamento de RH", Value: "empresa", Mapping: "organizacional"},
			{Label: "Escola ou universidade", Value: "escola", Mapping: "escolar"},
			{Label: "CRAS, CREAS ou ONG comunitaria", Value: "cras_creas", Mapping: "social"},
			{Label: "Hospital ou posto de saude", Value: "hospital", Mapping: "saude"},
			{Label: "Forum, penitenciaria ou instituicao juridica", Value: "forum", Mapping: "juridica"},
			{Label: "Clube esportivo ou centro de treinamento", Value: "clube_esportivo", Mapping: "esporte"},
			{Label: "Clinica de reabilitacao ou centro de neurociencias", Value: "clinica_reabilitacao", Mapping: "neuropsicologia"},
		},
	},
	{
		ID:    13,
		Text:  "Qual tipo de intervencao mais te interessa?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Psicoterapia individual ou de grupo", Value: "psicoterapia", Mapping: "clinica"},
			{Label: "Recrutamento, selecao e treinamento corporativo", Value: "recrutamento", Mapping: "organizacional"},
			{Label: "Orientacao educacional e apoio ao desenvolvimento de alunos", Value: "orientacao_educacional", Mapping: "escolar"},
			{Label: "Mobilizacao comunitaria e politicas publicas", Value: "mobilizacao", Mapping: "social"},
			{Label: "Acompanhamento psicologico hospitalar e cuidados paliativos", Value: "acompanhamento_hospitalar", Mapping: "saude"},
			{Label: "Pericias psicologicas, laudos e avaliacao judicial", Value: "pericias", Mapping: "juridica"},
			{Label: "Preparacao mental e psicologia da performance esportiva", Value: "preparacao_mental", Mapping: "esporte"},
			{Label: "Avaliacao neuropsicologica e reabilitacao cognitiva", Value: "avaliacao_neuro", Mapping: "neuropsicologia"},
		},
	},
	{
		ID:    14,
		Text:  "Descreva brevemente uma situacao profissional que voce se imagina vivendo como psicologo(a).",
		Type:  "open_ended",
		Block: "fields",
	},
	{
		ID:    15,
		Text:  "O que mais te motiva a ser psicologo(a)?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Aliviar o sofrimento emocional das pessoas", Value: "aliviar_sofrimento", Mapping: "clinica"},
			{Label: "Melhorar o desempenho e bem-estar nas organizacoes", Value: "desempenho_organizacoes", Mapping: "organizacional"},
			{Label: "Transformar a educacao e apoiar o desenvolvimento de jovens", Value: "transformar_educacao", Mapping: "escolar"},
			{Label: "Promover justica social e empoderar comunidades", Value: "justica_social", Mapping: "social"},
			{Label: "Cuidar da saude integral de pacientes em contexto medico", Value: "saude_integral", Mapping: "saude"},
			{Label: "Garantir direitos e contribuir para a justica", Value: "garantir_direitos", Mapping: "juridica"},
			{Label: "Potencializar a performance e o foco mental de atletas", Value: "potencializar_performance", Mapping: "esporte"},
			{Label: "Entender o cerebro e ajudar na reabilitacao cognitiva", Value: "entender_cerebro", Mapping: "neuropsicologia"},
		},
	},
}
