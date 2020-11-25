import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { EfComponent } from './ef.component';

describe('EfComponent', () => {
  let component: EfComponent;
  let fixture: ComponentFixture<EfComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ EfComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(EfComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
