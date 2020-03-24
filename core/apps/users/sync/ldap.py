import json

from django.contrib.auth.models import User
from ldap3 import Server, Connection

from kubeops_api.models.setting import Setting


class LDAPSync:
    def __init__(self):
        self._conn = None
        settings = Setting.get_settings(tab='ldap')
        self.ldap_enable = settings.get("AUTH_LDAP_ENABLE", False)
        if not self.ldap_enable:
            return
        self.bind_dn = settings.get("AUTH_LDAP_BIND_DN")
        self.bind_password = settings.get("AUTH_LDAP_BIND_PASSWORD")
        self.search_ou = settings.get("AUTH_LDAP_SEARCH_OU")
        self.search_filter = settings.get("AUTH_LDAP_SEARCH_FILTER")
        self.server_uri = settings.get("AUTH_LDAP_SERVER_URI")
        self.attr_map = json.loads(settings.get("AUTH_LDAP_USER_ATTR_MAP"))

    @property
    def connection(self):
        if self._conn:
            return self._conn
        server = Server(self.server_uri, use_ssl=False)
        conn = Connection(server, self.bind_dn, self.bind_password)
        conn.bind()
        self._conn = conn
        return self._conn

    def search_users(self):
        user_entries = list()
        search_ous = str(self.search_ou).split('|')
        for ou in search_ous:
            self.search_user_entries_ou(search_ou=ou)
            user_entries.extend(self.connection.entries)
        return user_entries

    def search_user_entries_ou(self, search_ou):
        search_filter = self.search_filter % {'user': '*'}
        attributes = list(self.attr_map.values())
        self.connection.search(
            search_base=search_ou, search_filter=search_filter,
            attributes=attributes)

    def user_entry_to_dict(self, entry):
        user = {}
        attr_map = self.attr_map.items()
        for attr, mapping in attr_map:
            if not hasattr(entry, mapping):
                continue
            value = getattr(entry, mapping).value or ''
            user[attr] = value
        return user

    def user_entries_to_dict(self, user_entries):
        users = []
        for user_entry in user_entries:
            user = self.user_entry_to_dict(user_entry)
            users.append(user)
        return users

    def run(self):
        user_entries = self.search_users()
        user_dicts = self.user_entries_to_dict(user_entries)
        for ud in user_dicts:
            defaults = {
                "username": ud.get("username", None),
                "email": ud.get("email", None)
            }
            if not defaults["username"] or not defaults["email"]:
                continue
            User.objects.get_or_create(defaults, username=defaults.get("username"))
