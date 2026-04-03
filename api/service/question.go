package service

type Option struct {
	Label   string `json:"label"`
	Value   string `json:"value"`
	Mapping string `json:"mapping,omitempty"`
}

type Question struct {
	ID      int      `json:"id"`
	Text    string   `json:"text"`
	Type    string   `json:"type"`  // multiple_choice | open_ended
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
	// ===== BLOCO 1: ABORDAGENS TEÓRICAS (10 perguntas) =====
	{
		ID:    1,
		Text:  "Na sua opinião, qual é a principal causa do sofrimento humano?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Conflitos inconscientes não resolvidos", Value: "conflitos_inconscientes", Mapping: "psicanalise"},
			{Label: "A perda de sentido e a dificuldade de ser autêntico", Value: "perda_sentido", Mapping: "fenomenologia"},
			{Label: "Padrões de comportamento aprendidos que geram consequências negativas", Value: "padroes_comportamento", Mapping: "comportamental"},
			{Label: "Pensamentos distorcidos sobre si mesmo e o mundo", Value: "pensamentos_distorcidos", Mapping: "tcc"},
			{Label: "O desequilíbrio entre aspectos conscientes e inconscientes da personalidade", Value: "desequilibrio_personalidade", Mapping: "junguiana"},
			{Label: "A interrupção do contato consigo mesmo e com o ambiente", Value: "interrupcao_contato", Mapping: "gestalt"},
			{Label: "As condições sociais e históricas que oprimem o indivíduo", Value: "condicoes_sociais", Mapping: "socio_historica"},
			{Label: "Disfunções nos padrões de comunicação e relacionamento", Value: "disfuncoes_comunicacao", Mapping: "sistemica"},
		},
	},
	{
		ID:    2,
		Text:  "Qual visão sobre a natureza humana mais ressoa com você?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Somos movidos por desejos inconscientes e pulsões", Value: "desejos_inconscientes", Mapping: "psicanalise"},
			{Label: "Somos seres livres, responsáveis por criar nosso próprio sentido", Value: "seres_livres", Mapping: "fenomenologia"},
			{Label: "Somos produto do nosso ambiente e das consequências dos nossos atos", Value: "produto_ambiente", Mapping: "comportamental"},
			{Label: "Temos capacidade racional de reestruturar nossos pensamentos", Value: "capacidade_racional", Mapping: "tcc"},
			{Label: "Buscamos a individuação — integrar luz e sombra", Value: "individuacao", Mapping: "junguiana"},
			{Label: "Somos totalidades em constante processo de ajustamento criativo", Value: "totalidades", Mapping: "gestalt"},
			{Label: "Somos seres históricos, construídos nas relações sociais", Value: "seres_historicos", Mapping: "socio_historica"},
			{Label: "Somos parte de sistemas interconectados que nos influenciam", Value: "sistemas_interconectados", Mapping: "sistemica"},
		},
	},
	{
		ID:    3,
		Text:  "Em poucas palavras, o que significa 'cura' ou 'melhora' em psicologia para você?",
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
			{Label: "Estar presente e facilitar a exploração da experiência vivida", Value: "estar_presente", Mapping: "fenomenologia"},
			{Label: "Analisar e modificar comportamentos através de técnicas baseadas em evidências", Value: "analisar_comportamentos", Mapping: "comportamental"},
			{Label: "Ajudar a identificar e reestruturar pensamentos disfuncionais", Value: "reestruturar_pensamentos", Mapping: "tcc"},
			{Label: "Guiar o paciente na exploração de símbolos e arquétipos", Value: "explorar_simbolos", Mapping: "junguiana"},
			{Label: "Ampliar a awareness e o contato com o momento presente", Value: "ampliar_awareness", Mapping: "gestalt"},
			{Label: "Promover consciência crítica e transformação social", Value: "consciencia_critica", Mapping: "socio_historica"},
			{Label: "Entender e reorganizar os padrões relacionais do sistema", Value: "reorganizar_padroes", Mapping: "sistemica"},
		},
	},
	{
		ID:    5,
		Text:  "Qual dessas afirmações mais ressoa com você?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "\"O sonho é a estrada real para o inconsciente\"", Value: "sonho_inconsciente", Mapping: "psicanalise"},
			{Label: "\"Tornar-se pessoa é o objetivo da vida\"", Value: "tornar_pessoa", Mapping: "fenomenologia"},
			{Label: "\"O comportamento é função de suas consequências\"", Value: "comportamento_consequencias", Mapping: "comportamental"},
			{Label: "\"Não são as coisas que nos perturbam, mas a visão que temos delas\"", Value: "visao_coisas", Mapping: "tcc"},
			{Label: "\"Quem olha para fora sonha; quem olha para dentro desperta\"", Value: "olha_dentro", Mapping: "junguiana"},
			{Label: "\"Perder a cabeça e cair em si\"", Value: "cair_em_si", Mapping: "gestalt"},
			{Label: "\"O homem se faz ao fazer o mundo e o mundo o faz\"", Value: "homem_mundo", Mapping: "socio_historica"},
			{Label: "\"O problema não está na pessoa, está no padrão\"", Value: "problema_padrao", Mapping: "sistemica"},
		},
	},
	{
		ID:    6,
		Text:  "Se você pudesse usar uma única palavra para definir o objetivo da psicoterapia, qual seria?",
		Type:  "open_ended",
		Block: "approaches",
	},
	{
		ID:    7,
		Text:  "Que tipo de evidência você mais valoriza na psicologia?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "A escuta profunda do discurso e dos lapsos do paciente", Value: "escuta_profunda", Mapping: "psicanalise"},
			{Label: "A descrição detalhada da experiência vivida", Value: "experiencia_vivida", Mapping: "fenomenologia"},
			{Label: "Dados experimentais e observações sistemáticas do comportamento", Value: "dados_experimentais", Mapping: "comportamental"},
			{Label: "Estudos clínicos controlados e protocolos validados", Value: "estudos_controlados", Mapping: "tcc"},
			{Label: "A análise de sonhos, mitos e símbolos universais", Value: "sonhos_mitos", Mapping: "junguiana"},
			{Label: "A observação fenomenológica do aqui-e-agora", Value: "aqui_agora", Mapping: "gestalt"},
			{Label: "A análise do contexto social, político e histórico", Value: "contexto_social", Mapping: "socio_historica"},
			{Label: "A observação dos padrões de interação no sistema familiar", Value: "padroes_interacao", Mapping: "sistemica"},
		},
	},
	{
		ID:    8,
		Text:  "Qual dessas atividades profissionais mais te atrairia?",
		Type:  "multiple_choice",
		Block: "approaches",
		Options: []Option{
			{Label: "Conduzir sessões de análise com associação livre e interpretação de sonhos", Value: "analise_sonhos", Mapping: "psicanalise"},
			{Label: "Acompanhar alguém em uma crise existencial com escuta empática", Value: "crise_existencial", Mapping: "fenomenologia"},
			{Label: "Desenvolver um programa de modificação de comportamento", Value: "modificacao_comportamento", Mapping: "comportamental"},
			{Label: "Aplicar técnicas estruturadas para tratar ansiedade ou depressão", Value: "tecnicas_estruturadas", Mapping: "tcc"},
			{Label: "Facilitar um processo de autoconhecimento profundo usando símbolos", Value: "autoconhecimento_simbolos", Mapping: "junguiana"},
			{Label: "Conduzir um grupo de vivência focado na percepção corporal", Value: "vivencia_corporal", Mapping: "gestalt"},
			{Label: "Coordenar um projeto de intervenção comunitária", Value: "intervencao_comunitaria", Mapping: "socio_historica"},
			{Label: "Fazer terapia familiar trabalhando as dinâmicas do sistema", Value: "terapia_familiar", Mapping: "sistemica"},
		},
	},
	{
		ID:    9,
		Text:  "Qual filósofo, pensador ou autor mais influenciou (ou influenciaria) sua visão de mundo?",
		Type:  "open_ended",
		Block: "approaches",
	},
	{
		ID:    10,
		Text:  "Um paciente diz: 'Não sei mais quem eu sou.' Como você abordaria essa fala?",
		Type:  "open_ended",
		Block: "approaches",
	},
	// ===== BLOCO 2: CAMPOS DE ATUAÇÃO (5 perguntas) =====
	{
		ID:    11,
		Text:  "Com qual público você mais gostaria de trabalhar?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Crianças e adolescentes em sofrimento emocional", Value: "criancas_adolescentes", Mapping: "clinica"},
			{Label: "Adultos em crises pessoais e emocionais", Value: "adultos_crises", Mapping: "clinica"},
			{Label: "Atletas e esportistas", Value: "atletas", Mapping: "esporte"},
			{Label: "Funcionários e gestores de empresas", Value: "funcionarios_empresas", Mapping: "organizacional"},
			{Label: "Comunidades em situação de vulnerabilidade", Value: "comunidades", Mapping: "social"},
			{Label: "Pacientes hospitalizados ou em tratamento médico", Value: "pacientes_hospitalizados", Mapping: "saude"},
			{Label: "Pessoas envolvidas em processos judiciais ou em situação carcerária", Value: "processos_judiciais", Mapping: "juridica"},
			{Label: "Alunos e educadores em contexto escolar", Value: "alunos_educadores", Mapping: "escolar"},
		},
	},
	{
		ID:    12,
		Text:  "Qual ambiente de trabalho mais combina com você?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Consultório particular", Value: "consultorio", Mapping: "clinica"},
			{Label: "Empresa ou departamento de RH", Value: "empresa", Mapping: "organizacional"},
			{Label: "Escola ou universidade", Value: "escola", Mapping: "escolar"},
			{Label: "CRAS, CREAS ou ONG comunitária", Value: "cras_creas", Mapping: "social"},
			{Label: "Hospital ou posto de saúde", Value: "hospital", Mapping: "saude"},
			{Label: "Fórum, penitenciária ou instituição jurídica", Value: "forum", Mapping: "juridica"},
			{Label: "Clube esportivo ou centro de treinamento", Value: "clube_esportivo", Mapping: "esporte"},
			{Label: "Clínica de reabilitação ou centro de neurociências", Value: "clinica_reabilitacao", Mapping: "neuropsicologia"},
		},
	},
	{
		ID:    13,
		Text:  "Qual tipo de intervenção mais te interessa?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Psicoterapia individual ou de grupo", Value: "psicoterapia", Mapping: "clinica"},
			{Label: "Recrutamento, seleção e treinamento corporativo", Value: "recrutamento", Mapping: "organizacional"},
			{Label: "Orientação educacional e apoio ao desenvolvimento de alunos", Value: "orientacao_educacional", Mapping: "escolar"},
			{Label: "Mobilização comunitária e políticas públicas", Value: "mobilizacao", Mapping: "social"},
			{Label: "Acompanhamento psicológico hospitalar e cuidados paliativos", Value: "acompanhamento_hospitalar", Mapping: "saude"},
			{Label: "Perícias psicológicas, laudos e avaliação judicial", Value: "pericias", Mapping: "juridica"},
			{Label: "Preparação mental e psicologia da performance esportiva", Value: "preparacao_mental", Mapping: "esporte"},
			{Label: "Avaliação neuropsicológica e reabilitação cognitiva", Value: "avaliacao_neuro", Mapping: "neuropsicologia"},
		},
	},
	{
		ID:    14,
		Text:  "Descreva brevemente uma situação profissional que você se imagina vivendo como psicólogo(a).",
		Type:  "open_ended",
		Block: "fields",
	},
	{
		ID:    15,
		Text:  "O que mais te motiva a ser psicólogo(a)?",
		Type:  "multiple_choice",
		Block: "fields",
		Options: []Option{
			{Label: "Aliviar o sofrimento emocional das pessoas", Value: "aliviar_sofrimento", Mapping: "clinica"},
			{Label: "Melhorar o desempenho e bem-estar nas organizações", Value: "desempenho_organizacoes", Mapping: "organizacional"},
			{Label: "Transformar a educação e apoiar o desenvolvimento de jovens", Value: "transformar_educacao", Mapping: "escolar"},
			{Label: "Promover justiça social e empoderar comunidades", Value: "justica_social", Mapping: "social"},
			{Label: "Cuidar da saúde integral de pacientes em contexto médico", Value: "saude_integral", Mapping: "saude"},
			{Label: "Garantir direitos e contribuir para a justiça", Value: "garantir_direitos", Mapping: "juridica"},
			{Label: "Potencializar a performance e o foco mental de atletas", Value: "potencializar_performance", Mapping: "esporte"},
			{Label: "Entender o cérebro e ajudar na reabilitação cognitiva", Value: "entender_cerebro", Mapping: "neuropsicologia"},
		},
	},
}
