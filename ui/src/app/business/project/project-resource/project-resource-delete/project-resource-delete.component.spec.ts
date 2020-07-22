import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectResourceDeleteComponent } from './project-resource-delete.component';

describe('ProjectResourceDeleteComponent', () => {
  let component: ProjectResourceDeleteComponent;
  let fixture: ComponentFixture<ProjectResourceDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ProjectResourceDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectResourceDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
