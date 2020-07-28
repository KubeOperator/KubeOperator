import {Directive, Input, TemplateRef, ViewContainerRef} from '@angular/core';
import {PermissionService} from '../auth/permission.service';

@Directive({
    selector: '[appMenuAuth]'
})
export class MenuAuthDirective {

    constructor(
        private templateRef: TemplateRef<any>,
        private viewContainer: ViewContainerRef,
        private permission: PermissionService) {
    }

    @Input() set appMenuAuth(menu: string) {
        this.permission.authRootMenu(menu).then(value => {
            const auth = value;
            if (auth) {
                this.viewContainer.createEmbeddedView(this.templateRef);
            } else if (!auth) {
                this.viewContainer.clear();
            }
        });
    }

}
