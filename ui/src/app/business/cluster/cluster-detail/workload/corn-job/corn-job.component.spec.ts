import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CornJobComponent } from './corn-job.component';

describe('CornJobComponent', () => {
  let component: CornJobComponent;
  let fixture: ComponentFixture<CornJobComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CornJobComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CornJobComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
