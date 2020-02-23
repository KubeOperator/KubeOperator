from users.serializers import UserSerializer, ProfileSerializer


def jwt_response_payload_handler(token, user=None, request=None):
    return generate_profile(token, user, request)


def generate_profile(token, user=None, request=None) -> dict:
    profile = ProfileSerializer(user.profile, context={"request": request}).data
    token = {
        "token": token
    }
    profile["user"].update(token)
    return profile
