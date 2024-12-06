import React, { useState } from 'react';
import { View, StyleSheet } from 'react-native';
import { TextInput, Button, Text } from 'react-native-paper';
import { useMutation } from '@tanstack/react-query';
import { login } from '../../services/auth';

function LoginScreen({ navigation }) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const loginMutation = useMutation(login, {
    onSuccess: (data) => {
      // Store token and navigate to main app
      navigation.replace('Main');
    },
    onError: (error) => {
      // Handle error
      console.error('Login failed:', error);
    },
  });

  const handleLogin = () => {
    loginMutation.mutate({ email, password });
  };

  return (
    <View style={styles.container}>
      <Text style={styles.title} variant="headlineMedium">
        Welcome to CoinSight
      </Text>
      
      <TextInput
        label="Email"
        value={email}
        onChangeText={setEmail}
        mode="outlined"
        style={styles.input}
        keyboardType="email-address"
        autoCapitalize="none"
      />
      
      <TextInput
        label="Password"
        value={password}
        onChangeText={setPassword}
        mode="outlined"
        style={styles.input}
        secureTextEntry
      />
      
      <Button
        mode="contained"
        onPress={handleLogin}
        style={styles.button}
        loading={loginMutation.isLoading}>
        Login
      </Button>
      
      <Button
        mode="text"
        onPress={() => navigation.navigate('Register')}
        style={styles.button}>
        Don't have an account? Register
      </Button>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    padding: 20,
    justifyContent: 'center',
  },
  title: {
    textAlign: 'center',
    marginBottom: 30,
  },
  input: {
    marginBottom: 15,
  },
  button: {
    marginTop: 10,
  },
});

export default LoginScreen;
