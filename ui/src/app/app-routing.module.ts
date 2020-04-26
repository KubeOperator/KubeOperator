import {NgModule} from '@angular/core';
import {Routes, RouterModule} from '@angular/router';
import {AppComponent} from './app.component';
import {TestComponent} from './test/test.component';


const routes: Routes = [
    {
        path: '',
        component: AppComponent,
        children: [
            {path: 'test', component: TestComponent}
        ]

    },
];

@NgModule({
    imports: [RouterModule.forRoot(routes)],
    exports: [RouterModule]
})
export class AppRoutingModule {
}
