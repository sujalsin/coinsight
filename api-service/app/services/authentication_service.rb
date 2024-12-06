class AuthenticationService
  class AuthenticationError < StandardError; end
  
  def self.authenticate(email, password)
    user = User.find_by(email: email)
    
    unless user&.authenticate(password)
      raise AuthenticationError, 'Invalid email or password'
    end
    
    token = generate_token(user)
    { user: user, token: token }
  end
  
  def self.verify_token(token)
    begin
      decoded = JWT.decode(token, jwt_secret, true, algorithm: 'HS256')[0]
      User.find(decoded['user_id'])
    rescue JWT::DecodeError, ActiveRecord::RecordNotFound
      nil
    end
  end
  
  private
  
  def self.generate_token(user)
    payload = {
      user_id: user.id,
      exp: 24.hours.from_now.to_i
    }
    
    JWT.encode(payload, jwt_secret, 'HS256')
  end
  
  def self.jwt_secret
    ENV['JWT_SECRET'] || Rails.application.secrets.secret_key_base
  end
end
