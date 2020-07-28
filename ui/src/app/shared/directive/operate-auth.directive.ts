import {Directive, Input, TemplateRef, ViewContainerRef} from '@angular/core';
import {PermissionService} from '../auth/permission.service';

@Directive({
    selector: '[appOperateAuth]'
})
export class OperateAuthDirective {

    constructor(private templateRef: TemplateRef<any>,
                private viewContainer: ViewContainerRef,
                private permission: PermissionService) {
    }

    @Input() set appOperateAuth(op: {}) {
        const operate = op['op'];
        const projectId = op['resourceId'];
        this.permission.authOperate(operate, projectId).then(value => {
            const auth = value;
            if (auth) {
                this.viewContainer.createEmbeddedView(this.templateRef);
            } else if (!auth) {
                this.viewContainer.clear();
            }
        });
    }
}
