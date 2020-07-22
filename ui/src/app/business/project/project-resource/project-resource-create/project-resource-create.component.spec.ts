import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectResourceCreateComponent } from './project-resource-create.component';

describe('ProjectResourceCreateComponent', () => {
  let component: ProjectResourceCreateComponent;
  let fixture: ComponentFixture<ProjectResourceCreateComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectResourceCreateComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectResourceCreateComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
