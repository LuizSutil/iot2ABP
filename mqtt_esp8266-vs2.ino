/*
 Basic ESP8266 MQTT example
 This sketch demonstrates the capabilities of the pubsub library in combination
 with the ESP8266 board/library.
 It connects to an MQTT server then:
  - publishes "hello world" to the topic "outTopic" every two seconds
  - subscribes to the topic "inTopic", printing out any messages
    it receives. NB - it assumes the received payloads are strings not binary
  - If the first character of the topic "inTopic" is an 1, switch ON the ESP Led,
    else switch it off
 It will reconnect to the server if the connection is lost using a blocking
 reconnect function. See the 'mqtt_reconnect_nonblocking' example for how to
 achieve the same result without blocking the main loop.
 To install the ESP8266 board, (using Arduino 1.6.4+):
  - Add the following 3rd party board manager under "File -> Preferences -> Additional Boards Manager URLs":
       http://arduino.esp8266.com/stable/package_esp8266com_index.json
  - Open the "Tools -> Board -> Board Manager" and click install for the ESP8266"
  - Select your ESP8266 in "Tools -> Board"
*/

#include <ESP8266WiFi.h>
#include <PubSubClient.h>

// Update these with values suitable for your network.

const char* ssid = "Rodox";
const char* password = "rodox321";
const char* mqtt_server = "broker.hivemq.com";

#define ID_MQTT  "ABP_IOT_EDDIE"     //id mqtt (para identificação de sessão)
                               //IMPORTANTE: este deve ser único no broker (ou seja, 
                               //            se um client MQTT tentar entrar com o mesmo 
                               //            id de outro já conectado ao broker, o broker 
                               //            irá fechar a conexão de um deles).

WiFiClient espClient;
PubSubClient client(espClient);
unsigned long lastMsg = 0;
#define MSG_BUFFER_SIZE	(50)
char msg[MSG_BUFFER_SIZE];
int value = 0;
char Sala;  
char Quarto;                                                                                      
char Cozinha;  

//Define os pinos para o trigger e echo Sensor de presença
#define pino_trigger D13
#define pino_echo D12

float calcularDistancia();

void setup() {
  pinMode(BUILTIN_LED, OUTPUT);     // Initialize the BUILTIN_LED pin as an output
  Serial.begin(115200);

  Sala = D8;
  Quarto = D7;                                                                                     
  Cozinha = D4;                                                                                                                                                                        

  pinMode(Sala, OUTPUT);    
  pinMode(Quarto, OUTPUT);                                                                        
  pinMode(Cozinha, OUTPUT);                                                                      
 
  pinMode(pino_echo, INPUT);
  pinMode(pino_trigger, OUTPUT);
  
  setup_wifi();
  client.setServer(mqtt_server, 1883);
  client.setCallback(callback);
}


void setup_wifi() {

  delay(10);
  // We start by connecting to a WiFi network
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);

  WiFi.mode(WIFI_STA);
  WiFi.begin(ssid, password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  randomSeed(micros());

  Serial.println("");
  Serial.println("WiFi connected");
  Serial.println("IP address: ");
  Serial.println(WiFi.localIP());
}

void callback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();

   /*
  // Switch on the LED if an 1 was received as first character
  if ((char)payload[0] == '1') {
    digitalWrite(BUILTIN_LED, LOW);   // Turn the LED on (Note that LOW is the voltage level
    // but actually the LED is on; this is because
    // it is active low on the ESP-01)
  } else {
    digitalWrite(BUILTIN_LED, HIGH);  // Turn the LED off by making the voltage HIGH
  }*/

  Serial.println("");                                                                             
  if ((char)payload[0] == 'S') {                                                                  
    digitalWrite(Sala, HIGH);                                                                     
    //snprintf (msg, MSG_BUFFER_SIZE, "A luz da sala está ligada");                                 
    //Serial.print("Publica mensagem: ");                                                           
    //Serial.println(msg);                                                                          
    //client.publish("casa/sala", msg);                                                             
  }
  Serial.println("");                                                                             
  if ((char)payload[0] == 's') {                                                                  
    digitalWrite(Sala, LOW);                                                                      
    //snprintf (msg, MSG_BUFFER_SIZE, "A luz da sala está desligada");                              
    //Serial.print("Publica mensagem: ");                                                           
    //Serial.println(msg);                                                                          
    //client.publish("casa/sala", msg);                                                             
  }
  Serial.println("");  
  if ((char)payload[0] == 'Q') {                                                                  
    digitalWrite(Quarto, HIGH);                                                                   
    //snprintf (msg, MSG_BUFFER_SIZE, "A luz do quarto está ligada");                               
    //Serial.print("Publica mensagem: ");                                                           
    //Serial.println(msg);                                                                          
    //client.publish("casa/quarto", msg);                                                           
  }
   Serial.println("");                                                                            
   if ((char)payload[0] == 'q') {                                                                 
    digitalWrite(Quarto, LOW);                                                                    
    //snprintf (msg, MSG_BUFFER_SIZE, "A luz do quarto está desligada");                            
    //Serial.print("Publica mensagem: ");                                                           
    //Serial.println(msg);                                                                          
    // client.publish("casa/quarto", msg);                                                           
  }
  Serial.println("");                                                                             
  if ((char)payload[0] == 'C') {                                                                  
    digitalWrite(Cozinha, HIGH);                                                                  
//    snprintf (msg, MSG_BUFFER_SIZE, "A luz da cozinha está ligada");                              
//    Serial.print("Publica mensagem: ");                                                           
//    Serial.println(msg);                                                                          
//    client.publish("casa/cozinha", msg);                                                          
  }
  Serial.println("");                                                                             
   if ((char)payload[0] == 'c') {                                                                 
    digitalWrite(Cozinha, LOW);                                                                   
//    snprintf (msg, MSG_BUFFER_SIZE, "A luz da cozinha está desligada");                           
//    Serial.print("Publica mensagem: ");                                                           
//    Serial.println(msg);                                                                          
//    client.publish("satc/iot/luiz/cozinha", msg);                                                          
  }
  Serial.println(""); 

}

void reconnect() {
  // Loop until we're reconnected
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");

    /*
    // Create a random client ID
    String clientId = "ESP8266Client-";
    clientId += String(random(0xffff), HEX);
    client.connect(clientId.c_str())
    */
    
    // Attempt to connect
    if (client.connect(ID_MQTT)) {
      Serial.println("connected");
      // Once connected, publish an announcement...
      //client.publish("satc/iot/eddie", "hello world");
      // ... and resubscribe
      client.subscribe("satc/iot/eddie");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      // Wait 5 seconds before retrying
      delay(5000);
    }
  }
}



void loop() {

  if (!client.connected()) {
    reconnect();
  }
  client.loop();

  unsigned long now = millis();
  if (now - lastMsg > 2000) {
    lastMsg = now;
    ++value;
    snprintf (msg, MSG_BUFFER_SIZE, "hello world #%ld", value);
    Serial.print("Publish message: ");
    Serial.println(msg);

    char buffer[80];
    char fValue[16];
    dtostrf(calcularDistancia(), 3, 2, fValue);
    snprintf(buffer, sizeof(buffer), "%s", fValue);

    client.publish("satc/iot/eddie", buffer);
    
    //Serial.print(distanciacm);
    //client.publish("satc/iot/eddie", msg);
    //client.publish("satc/iot/eddie/sensor",  );
  }
}

float calcularDistancia(){
  digitalWrite(pino_trigger,HIGH);
  delay(500);
  digitalWrite(pino_trigger,LOW);
  unsigned long time = pulseIn(pino_echo, HIGH);
  float distancia = time*0.0171;

  //Serial.print(distancia);
  //Serial.print(" cm");
  //Serial.println();
  
  //Serial.println("\n" + String(distancia) + ");
  return distancia;
}
