# Questão 3

Considere o seguinte estudo de caso:

Uma plataforma de notícias, chamada absurd-news.net, provê aos seus usuários notícias sobre os mais variados assuntos 24 horas por dia. Para acessar de forma ilimitada a plataforma, o usuário necessita se cadastrar na plataforma e pagar um valor mensal não especificado.
Os benefícios de assinar ao site são inúmeros, dentre os quais a possibilidade de acessar todo o seu histórico de noticias visualizadas nos últimos 6 meses. A plataforma recebe milhões de acessos por dia, considerando que é um dos sites mais acessados de seu continente de origem.
No entanto, a funcionalidade de visualização do histórico de noticias se tornou incrivelmente não responsiva nos últimos meses. Essa funcionalidade é provida por uma tabela relacional simples, chamada `history`, que contém os metadados da notícia, como o identificador único da notícia, assim como, o identificador único do usuário. A tabela contém uma chave única `(user_id, news_id)` que é indexada. A tabela contém mais de 200 milhões de registros, e mesmo uma busca simples de um histórico de um usuário leva mais de 1s para completar.
Com este cenário, o número de usuários da plataforma começou a cair, considerando que a funcionalidade de histórico, muito utilizada por mais de 90% dos usuários, não é responsiva em termos de latência.
Com este cenário em mente, a plataforma deseja que 2 requisitos não funcionais sejam atingidos:

- A funcionalidade de histórico deve funcionar mesmo que o banco de dados esteja
  completamente offline.
- Nenhuma visita pode ser perdida no cenário de indisponibilidade dos sistemas de
  armazenamento.
- O tempo de latência para qualquer requisição de busca de histórico deve ser inferior a 10ms
  em 80% dos casos.

Sugira as modificações necessárias para atingir estes objetivos

## Solução

Para resolver esse problema e atender aos requisitos específicos, pode-se utilizar a seguinte abordagem.

Primeiramente, é recomendável implementar a paginação para limitar o número de notícias retornadas em cada consulta de histórico. Além disso, técnicas de otimização, como a criação de índices adicionais, podem ser aplicadas para acelerar as buscas na tabela `history`.

Uma abordagem viável consiste em possibilitar que uma porção do histórico de notícias visualizadas seja armazenada de forma local no navegador do cliente. Essa estratégia permite a atualização assíncrona desses dados com o banco de dados principal, assegurando a consistência mesmo em cenários de conectividade intermitente ou breves períodos de indisponibilidade do banco de dados central.

Para melhorar significativamente o tempo de resposta, é crucial implementar um sistema de cache robusto. Esse sistema pode armazenar resultados de consultas frequentes, reduzindo assim o tempo de latência para as requisições de busca de histórico.

Por fim, introduzir um mecanismo de Fallback, permitindo que os usuários acessem notícias mesmo quando o banco de dados principal estiver offline. Isso pode ser realizado através de uma camada de armazenamento local no servidor, proporcionando acesso limitado às notícias durante períodos de indisponibilidade.
