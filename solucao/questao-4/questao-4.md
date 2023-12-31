# Questão 4

Considere o seguinte estudo de caso:

Uma empresa, chamada X Company, produz uma solução para geração e envio de boletos de empresas prestadoras de serviços aos seus clientes finais. A aplicação funciona de maneira simples, a empresa se registra na plataforma, cadastra seus clientes, com suas respectivas informações (CPF, RG, nome, etc), e detalha o valor e periodicidade das cobranças. O sistema gera os boletos dentro do período de cobrança e os envia aos devedores para pagamento. Dado a natureza da aplicação, o usuário final não precisa se cadastrar para ter acesso aos boletos, ele deverá acessar o documento por meio de um link enviado no e-mail citado anteriormente. É de desejo da X Company que o usuário continue tendo acesso aos boletos sem se cadastrar na plataforma. O link enviado por email corresponde ao seguinte endpoint:

```plaintext
GET https://boletos.x-company.com/boletos/{id}
```

Desta forma, o cliente abre o link anexado ao email e acessa o boleto para pagamento, este contendo algumas informações sensíveis do usuário, como CPF e endereço Portanto, com base no cenário e aplicação proposta, identifique a existência (ou não) de vulnerabilidades de segurança e como, caso existam, estas podem ser usadas para gerar danos. Pontos extras para propostas de solução válidas das vulnerabilidades, caso existam. Para esta análise, considere algumas requisições logadas por meio do monitoramento de rede da empresa.

```plaintext
GET https://boletos.x-company.com/boletos/13857 
GET https://boletos.x-company.com/boletos/913851365 
GET https://boletos.x-company.com/boletos/1359861
```

## Solução  

Através da atual forma de acesso do cliente ao boleto, identificamos potenciais vulnerabilidades que podem ser exploradas por indivíduos mal-intencionados. Uma abordagem comum para aproveitar essas brechas de segurança é a utilização de scripts que geram números inteiros incrementalmente, realizando requisições no endpoint `GET https://boletos.x-company.com/boletos/{id}` até obter uma resposta positiva. Isso permite o acesso não autorizado aos dados dos clientes das empresas contratantes, resultando em danos significativos e preocupações para os usuários afetados pelo vazamento de informações.

Uma solução para isso, é a criação de IDs criptografados. Inibindo as tentativas de acesso as informações sensíveis, essas criptografias podem ser feitas através da geração de UUIDs  (Identificadores Únicos Universais) para cada boleto criado. Esses identificadores criptografados dificultam as tentativas de acesso não autorizado, uma vez que não é possível prever ou gerar sistematicamente os IDs válidos.

Além disso, é crucial complementar essa solução com a implementação de logs e monitoramentos proativos. Estabelecer sistemas que identifiquem e respondam a atividades suspeitas, como tentativas repetidas de acesso, é essencial para detectar possíveis ataques antes que causem danos substanciais.

Outra medida preventiva recomendada é a introdução de uma política de expiração para os boletos. Ao atribuir uma validade limitada aos links enviados aos clientes, especialmente após a confirmação do pagamento, reduz-se significativamente a janela de oportunidade para possíveis ataques. Isso garante que os boletos se tornem inacessíveis após um período determinado, aumentando a segurança e protegendo as informações sensíveis dos clientes.
