# pylint: skip-file
# flake8: noqa


# pylint: disable=too-many-instance-attributes
class OCLabel(OpenShiftCLI):
    ''' Class to wrap the oc command line tools '''

    # pylint allows 5
    # pylint: disable=too-many-arguments
    def __init__(self,
                 name,
                 namespace,
                 kind,
                 kubeconfig,
                 labels=None,
                 selector=None,
                 verbose=False):
        ''' Constructor for OCLabel '''
        super(OCLabel, self).__init__(namespace, kubeconfig=kubeconfig, verbose=verbose)
        self.name = name
        self.kind = kind
        self.labels = labels
        self._curr_labels = None
        self.selector = selector

    @property
    def current_labels(self):
        '''property for the current labels'''
        if self._curr_labels is None:
            results = self.get()
            self._curr_labels = results['labels']

        return self._curr_labels

    @current_labels.setter
    def current_labels(self, data):
        '''property setter for current labels'''
        self._curr_labels = data

    def compare_labels(self, host_labels):
        ''' compare incoming labels against current labels'''

        for label in self.labels:
            if label['key'] not in host_labels or \
               label['value'] != host_labels[label['key']]:
                return False
        return True

    def all_user_labels_exist(self):
        ''' return whether all the labels already exist '''

        for current_host_labels in self.current_labels:
            rbool = self.compare_labels(current_host_labels)
            if not rbool:
                return False
        return True

    def any_label_exists(self):
        ''' return whether any single label already exists '''

        for current_host_labels in self.current_labels:
            for label in self.labels:
                if label['key'] in current_host_labels:
                    return True
        return False

    def get_user_keys(self):
        ''' go through list of user key:values and return all keys '''

        user_keys = []
        for label in self.labels:
            user_keys.append(label['key'])

        return user_keys

    def get_current_label_keys(self):
        ''' collect all the current label keys '''

        current_label_keys = []
        for current_host_labels in self.current_labels:
            for key in current_host_labels.keys():
                current_label_keys.append(key)

        return list(set(current_label_keys))

    def get_extra_current_labels(self):
        ''' return list of labels that are currently stored, but aren't
            in user-provided list '''

        extra_labels = []
        user_label_keys = self.get_user_keys()
        current_label_keys = self.get_current_label_keys()

        for current_key in current_label_keys:
            if current_key not in user_label_keys:
                extra_labels.append(current_key)

        return extra_labels

    def extra_current_labels(self):
        ''' return whether there are labels currently stored that user
            hasn't directly provided '''
        extra_labels = self.get_extra_current_labels()

        if len(extra_labels) > 0:
            return True

        return False

    def replace(self):
        ''' replace currently stored labels with user provided labels '''
        cmd = self.cmd_template()

        # First delete any extra labels
        extra_labels = self.get_extra_current_labels()
        if len(extra_labels) > 0:
            for label in extra_labels:
                cmd.append("{}-".format(label))

        # Now add/modify the user-provided label list
        if len(self.labels) > 0:
            for label in self.labels:
                cmd.append("{}={}".format(label['key'], label['value']))

        # --overwrite for the case where we are updating existing labels
        cmd.append("--overwrite")
        return self.openshift_cmd(cmd)

    def get(self):
        '''return label information '''

        result_dict = {}
        label_list = []

        if self.name:
            result = self._get(resource=self.kind, name=self.name, selector=self.selector)

            if result['results'][0] and 'labels' in result['results'][0]['metadata']:
                label_list.append(result['results'][0]['metadata']['labels'])
            else:
                label_list.append({})

        else:
            result = self._get(resource=self.kind, selector=self.selector)

            for item in result['results'][0]['items']:
                if 'labels' in item['metadata']:
                    label_list.append(item['metadata']['labels'])
                else:
                    label_list.append({})

        self.current_labels = label_list
        result_dict['labels'] = self.current_labels
        result_dict['item_count'] = len(self.current_labels)
        result['results'] = result_dict

        return result

    def cmd_template(self):
        ''' boilerplate oc command for modifying lables on this object '''
        # let's build the cmd with what we have passed in
        cmd = ["label", self.kind]

        if self.selector:
            cmd.extend(["--selector", self.selector])
        elif self.name:
            cmd.extend([self.name])

        return cmd

    def add(self):
        ''' add labels '''
        cmd = self.cmd_template()

        for label in self.labels:
            cmd.append("{}={}".format(label['key'], label['value']))

        cmd.append("--overwrite")

        return self.openshift_cmd(cmd)

    def delete(self):
        '''delete the labels'''
        cmd = self.cmd_template()
        for label in self.labels:
            cmd.append("{}-".format(label['key']))

        return self.openshift_cmd(cmd)

    # pylint: disable=too-many-branches,too-many-return-statements
    @staticmethod
    def run_ansible(params, check_mode=False):
        ''' run the oc_label module

            prams comes from the ansible portion of this module
            check_mode: does the module support check mode. (module.check_mode)
        '''
        oc_label = OCLabel(params['name'],
                           params['namespace'],
                           params['kind'],
                           params['kubeconfig'],
                           params['labels'],
                           params['selector'],
                           verbose=params['debug'])

        state = params['state']
        name = params['name']
        selector = params['selector']

        api_rval = oc_label.get()

        #####
        # Get
        #####
        if state == 'list':
            return {'changed': False, 'results': api_rval['results'], 'state': "list"}

        #######
        # Add
        #######
        if state == 'add':
            if not (name or selector):
                return {'failed': True,
                        'msg': "Param 'name' or 'selector' is required if state == 'add'"}
            if not oc_label.all_user_labels_exist():
                if check_mode:
                    return {'changed': False, 'msg': 'Would have performed an addition.'}
                api_rval = oc_label.add()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "add"}

            return {'changed': False, 'state': "add"}

        ########
        # Delete
        ########
        if state == 'absent':
            if not (name or selector):
                return {'failed': True,
                        'msg': "Param 'name' or 'selector' is required if state == 'absent'"}

            if oc_label.any_label_exists():
                if check_mode:
                    return {'changed': False, 'msg': 'Would have performed a delete.'}

                api_rval = oc_label.delete()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "absent"}

            return {'changed': False, 'state': "absent"}

        if state == 'present':
            ########
            # Update
            ########
            if not (name or selector):
                return {'failed': True,
                        'msg': "Param 'name' or 'selector' is required if state == 'present'"}
            # if all the labels passed in don't already exist
            # or if there are currently stored labels that haven't
            # been passed in
            if not oc_label.all_user_labels_exist() or \
               oc_label.extra_current_labels():
                if check_mode:
                    return {'changed': False, 'msg': 'Would have made changes.'}

                api_rval = oc_label.replace()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                # return the created object
                api_rval = oc_label.get()

                if api_rval['returncode'] != 0:
                    return {'failed': True, 'msg': api_rval}

                return {'changed': True, 'results': api_rval, 'state': "present"}

            return {'changed': False, 'results': api_rval, 'state': "present"}

        return {'failed': True,
                'changed': False,
                'results': 'Unknown state passed. %s' % state,
                'state': "unknown"}
