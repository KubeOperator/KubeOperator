import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { ShellComponent } from './base/shell/shell.component';
import { ProjectComponent } from './project/project.component';
import { ProjectDetailComponent } from './project/project-detail/project-detail.component';
import { ProjectOverviewComponent } from './project/project-overview/project-overview.component';
import { PlaybookComponent } from './playbook/playbook.component';
import { PlaybookCreateComponent } from './playbook/playbook-create/playbook-create.component';
import { RoleListComponent } from './role/role-list/role-list.component';
import { AdhocComponent } from './adhoc/adhoc/adhoc.component';
import { HostComponent } from './host/host.component';
import { GroupComponent } from './group/group.component';
import { PageNotFoundComponent } from './shared/not-found/not-found.component';

import { AuthCheckGuard } from './shared/route/auth-user-active.service';
import { ProjectRoutingResolver } from './project/project-routing-resolver.service';


export const ROUTES: Routes = [
  {path: '', redirectTo: 'ansible', pathMatch: 'full'},
  {
    path: 'ansible',
    component: ShellComponent,
    canActivateChild: [AuthCheckGuard],
    children: [
      {path: '', redirectTo: 'projects', pathMatch: 'full'},
      {
        path: 'projects/:project',
        component: ProjectDetailComponent,
        resolve: { project: ProjectRoutingResolver},
        canActivate: [AuthCheckGuard],
        canActivateChild: [AuthCheckGuard],
        children: [
          {path: '', redirectTo: 'overview', pathMatch: 'full'},
          {path: 'overview', component: ProjectOverviewComponent},
          {path: 'playbooks/create', component: PlaybookCreateComponent},
          {path: 'playbooks', component: PlaybookComponent},
          {path: 'roles', component: RoleListComponent},
          {path: 'adhoc', component: AdhocComponent},
          {path: 'hosts', component: HostComponent},
          {path: 'groups', component: GroupComponent},
        ]
      },
      {
        path: 'projects',
        component: ProjectComponent,
        children: [
        ]
      },
    ]
  },
  {path: '**', component: PageNotFoundComponent}
];

@NgModule({
  imports: [
    RouterModule.forRoot(ROUTES),
  ],
  exports: [
    RouterModule
  ]

})
export class AppRoutingModule {

}
