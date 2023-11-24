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

Através da atual forma de acesso do cliente ao boleto, existem brechas de vulnerabilidade que podem ser exploradas por pessoas más intencionadas. É possível utilizar um script para gerar números inteiros sendo incrementados em 1 e tentar fazer uma requisição no endpoint `GET https://boletos.x-company.com/boletos/{id}` até conseguir uma resposta de requisição positiva, acessando desta forma os dados dos clientes das empresas contratantes. Gerando danos e preocupações para os usuários que tiveram as informações vazadas.

Uma solução para isso, é a criação de IDs criptografados para inibir às tentativas de acesso as informações sensíveis, essas criptografias podem ser feitas através da geração de UUIDs para cada boleto criado. Complementar a essa solução, seria de extrema importância implementar logs e monitoramentos para identificar e responder a atividades suspeitas, como por exemplo, tentativas repetidas de acesso.

Outra solução que poderia ser adotada é implementar uma expiração para o boleto assim que constar como pago. Fazendo com que os links enviados aos clientes tenham uma validade e reduzindo a janela de oportunidades para possíveis ataques.
