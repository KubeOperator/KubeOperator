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
        const projectName = op['projectName'];
        this.permission.authOperate(operate, projectName).then(value => {
            const auth = value;
            if (auth) {
                this.viewContainer.createEmbeddedView(this.templateRef);
            } else if (!auth) {
                this.viewContainer.clear();
            }
        });
    }
}
