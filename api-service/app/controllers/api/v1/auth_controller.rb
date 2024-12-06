module Api
  module V1
    class AuthController < ApplicationController
      skip_before_action :authenticate_user, only: [:login, :register]
      
      def login
        result = AuthenticationService.authenticate(
          params[:email],
          params[:password]
        )
        
        render json: {
          token: result[:token],
          user: UserSerializer.new(result[:user]).as_json
        }
      rescue AuthenticationService::AuthenticationError => e
        render json: { error: e.message }, status: :unauthorized
      end
      
      def register
        user = User.new(user_params)
        
        if user.save
          result = AuthenticationService.authenticate(
            user.email,
            params[:password]
          )
          
          render json: {
            token: result[:token],
            user: UserSerializer.new(user).as_json
          }, status: :created
        else
          render json: { errors: user.errors.full_messages },
                 status: :unprocessable_entity
        end
      end
      
      def me
        render json: UserSerializer.new(current_user).as_json
      end
      
      private
      
      def user_params
        params.permit(:email, :password, :username)
      end
    end
  end
end
